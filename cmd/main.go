package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	logSystem "github.com/renderview-inc/backend/internal/app/application/services/logger"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/option"
	"github.com/renderview-inc/backend/pkg/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/renderview-inc/backend/internal/app/application/middleware"
	"github.com/renderview-inc/backend/internal/app/application/services"
	"github.com/renderview-inc/backend/internal/app/infrastructure/cache"
	"github.com/renderview-inc/backend/internal/app/infrastructure/repositories"
	v1 "github.com/renderview-inc/backend/internal/app/presentation/api/handlers/v1"
	"github.com/renderview-inc/backend/internal/pkg/txhelper"
	postgres "github.com/renderview-inc/backend/pkg/connections"
)

func main() {
	ctx := context.Background()

	logService, err := registerLogService()
	if err != nil {
		log.Fatalf("failed to initialize log service: %v", err)
	}
	defer func(logService *logSystem.LogService) {
		err = logService.Sync()
		if err != nil {
			log.Fatalf("failed to sync log service: %v", err)
		}
	}(logService)

	dbPool, err := connectPostgres()
	if err != nil {
		logService.Error(ctx, "unable to connect to database", option.Error(err))

		return
	}
	defer dbPool.Close()

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	httpServerAddr := os.Getenv("HTTP_ADDR")

	userAccountRepo := repositories.NewUserAccountRepository(dbPool)
	userSessionRepo := repositories.NewUserSessionRepository(dbPool)
	loginHistoryRepo := repositories.NewLoginHistoryRepository(dbPool)
	chatRepo := repositories.NewChatRepository(dbPool)
	messageRepo := repositories.NewMessageRepository(dbPool)

	sessionCache := cache.NewUserSessionCache(redisAddr, redisPassword, 0)

	passwordHasher := services.NewBcryptPasswordHasher()
	txHelper := txhelper.NewTxHelper(dbPool)
	tokenIssuer := services.NewBase64TokenIssuer(20, 30*time.Minute, 30*24*time.Hour)
	tokenHasher := services.NewSha256TokenHasher()

	userAccountService := services.NewUserAccountService(userAccountRepo, passwordHasher)
	authService := services.NewAuthService(
		loginHistoryRepo,
		userSessionRepo,
		sessionCache,
		userAccountService,
		txHelper,
		tokenIssuer,
		tokenHasher,
	)
	chatService := services.NewChatService(chatRepo)
	messageService := services.NewMessageService(messageRepo)

	userAccountHandler := v1.NewUserAccountHandler(&userAccountService, passwordHasher)
	authHandler := v1.NewAuthHandler(authService)
	chatHandler := v1.NewChatHandler(chatService)
	messageHandler := v1.NewMessageHandler(messageService)

	r := mux.NewRouter()
	r.Use(middleware.CorrelationMiddleware)
	r.Use(func(next http.Handler) http.Handler {
		return middleware.LoggingMiddleware(next, logService)
	})

	public := r.NewRoute().Subrouter()
	protected := r.NewRoute().Subrouter()

	public.HandleFunc("/api/v1/user/register", userAccountHandler.HandleRegister).Methods(http.MethodPost)
	public.HandleFunc("/api/v1/auth/login", authHandler.HandleLogin).Methods(http.MethodPost)

	protected.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, authService)
	})

	protected.HandleFunc("/api/v1/auth/logout", authHandler.HandleLogout).Methods(http.MethodPost)
	protected.HandleFunc("/api/v1/auth/refresh", authHandler.HandleRefresh).Methods(http.MethodPost)

	protected.HandleFunc("/api/v1/chat", chatHandler.HandleCreateChat).Methods(http.MethodPost)
	protected.HandleFunc("/api/v1/chat/participant", chatHandler.HandleAddParticipant).Methods(http.MethodPost)
	protected.HandleFunc("/api/v1/chat/tag", chatHandler.HandleGetChatInfoByTag).Methods(http.MethodGet)
	protected.HandleFunc("/api/v1/chat/id", chatHandler.HandleGetChatInfoByID).Methods(http.MethodGet)
	protected.HandleFunc("/api/v1/chat", chatHandler.HandleGetChatsWithLastMessages).Methods(http.MethodGet)
	protected.HandleFunc("/api/v1/chat", chatHandler.HandleUpdateChat).Methods(http.MethodPut)
	protected.HandleFunc("/api/v1/chat", chatHandler.HandleDeleteChat).Methods(http.MethodDelete)
	protected.HandleFunc("/api/v1/chat/participant", chatHandler.HandleRemoveParticipant).Methods(http.MethodDelete)

	protected.HandleFunc("/api/v1/message", messageHandler.HandleCreateMessage).Methods(http.MethodPost)
	protected.HandleFunc("/api/v1/message", messageHandler.HandleGetMessage).Methods(http.MethodGet)
	protected.HandleFunc("/api/v1/message/last", messageHandler.HandleGetLastMessageByChatTag).Methods(http.MethodGet)
	protected.HandleFunc("/api/v1/message", messageHandler.HandleUpdateMessage).Methods(http.MethodPut)
	protected.HandleFunc("/api/v1/message", messageHandler.HandleDeleteMessage).Methods(http.MethodDelete)

	logService.Info(ctx, "starting server", option.Any("httpAddr", httpServerAddr))
	if err = http.ListenAndServe(httpServerAddr, r); err != nil {
		logService.Fatal(ctx, "failed to start server", option.Error(err))
	}
}

func connectPostgres() (*pgxpool.Pool, error) {
	dbURL := "postgres://" + os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	return postgres.NewPsqlPool(dbURL)
}

func registerLogService() (*logSystem.LogService, error) {
	cfg, err := config.LoadLogConfig()
	if err != nil {
		return nil, err
	}

	loggerService, err := logSystem.NewLogService(cfg)
	if err != nil {
		return nil, err
	}

	return loggerService, nil
}
