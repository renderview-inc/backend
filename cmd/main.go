package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/renderview-inc/backend/internal/app/application/services"
	"github.com/renderview-inc/backend/internal/app/infrastructure/cache"
	"github.com/renderview-inc/backend/internal/app/infrastructure/repositories"
	"github.com/renderview-inc/backend/internal/app/presentation/api/handlers/v1"
	"github.com/renderview-inc/backend/internal/pkg/txhelper"
	"github.com/renderview-inc/backend/pkg/connections"
)

func main() {
	dbURL := "postgres://" + os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	httpServerAddr := os.Getenv("HTTP_ADDR")

	// Подключение к PostgreSQL
	dbPool, err := postgres.NewPsqlPool(dbURL)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)

		return
	}
	defer dbPool.Close()

	// Репозитории
	userAccountRepo := repositories.NewUserAccountRepository(dbPool)
	userSessionRepo := repositories.NewUserSessionRepository(dbPool)
	loginHistoryRepo := repositories.NewLoginHistoryRepository(dbPool)

	// Redis
	sessionCache := cache.NewUserSessionCache(redisAddr, redisPassword, 0)

	// Сервисы и вспомогательные объекты
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

	// Handlers
	userAccountHandler := v1.NewUserAccountHandler(&userAccountService, passwordHasher)
	authHandler := v1.NewAuthHandler(authService)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/user/register", userAccountHandler.HandleRegister)
	mux.HandleFunc("/api/v1/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.HandleLogout)
	mux.HandleFunc("/api/v1/auth/refresh", authHandler.HandleRefresh)

	// Запуск сервера
	log.Printf("Starting server on %s\n", httpServerAddr)
	if err := http.ListenAndServe(httpServerAddr, mux); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
