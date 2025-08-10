package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/application/services"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type UserAccountHandler struct {
	accountService *services.UserAccountService
	passwordHasher services.BcryptPasswordHasher
	validate       *validator.Validate
}

func NewUserAccountHandler(accountService *services.UserAccountService, passwordHasher services.BcryptPasswordHasher) *UserAccountHandler {
	return &UserAccountHandler{
		accountService: accountService,
		passwordHasher: passwordHasher,
		validate:       validator.New(),
	}
}

func (uah *UserAccountHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var registerDto dtos.RegisterDto
	if err := json.NewDecoder(r.Body).Decode(&registerDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uah.validate.Struct(registerDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Manual validation for either email or phone
	if registerDto.Credentials.Email == "" && registerDto.Credentials.Phone == "" {
		http.Error(w, "Either email or phone must be provided", http.StatusBadRequest)
		return
	}

	hashedPassword, err := uah.passwordHasher.HashPassword(registerDto.Credentials.Password)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	userAccount := entities.NewUserAccount(
		uuid.New(),
		registerDto.Credentials.Tag,
		registerDto.Name,
		registerDto.Desc,
		hashedPassword,
		registerDto.Credentials.Email,
		registerDto.Credentials.Phone,
	)

	if err := uah.accountService.Register(r.Context(), userAccount); err != nil {
		// TODO: handle different error types
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
