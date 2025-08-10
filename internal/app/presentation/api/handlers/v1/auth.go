package v1

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/application/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&loginDto.Credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	loginDto.LoginMeta = dtos.LoginMetaDto{
		UserAgent: r.UserAgent(),
		IpAddr:    net.ParseIP(ip),
	}

	tokens, err := h.authService.Login(r.Context(), loginDto)
	if err != nil {
		// TODO: handle different error types
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var tokensDto dtos.TokensDto
	if err := json.NewDecoder(r.Body).Decode(&tokensDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.Logout(r.Context(), tokensDto); err != nil {
		// TODO: handle different error types
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var tokensDto dtos.TokensDto
	if err := json.NewDecoder(r.Body).Decode(&tokensDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), tokensDto.RefreshToken)
	if err != nil {
		// TODO: handle different error types
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
