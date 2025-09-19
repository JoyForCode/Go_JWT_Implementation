package handlers

import (
	"jwt_clean/internal/auth"
	"jwt_clean/internal/service"
	"jwt_clean/error"
	"net/http"
)

type AuthHandler struct {
	authService *service.Service
}

func NewAuthHandler(authService *service.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq auth.LoginRequest
	if err := apperror.ValidateJSON(r, &loginReq); err != nil {
		apperror.LogWarning("Invalid JSON in login request: %v", err)
		apperror.WriteError(w, http.StatusBadRequest, apperror.ErrInvalidJSON)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		apperror.WriteError(w, http.StatusBadRequest, apperror.ErrMissingCredentials)
		return
	}

	tokenPair, err := h.authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		apperror.LogWarning("Login failed for user: %s, error: %v", loginReq.Username, err)
		apperror.WriteError(w, http.StatusUnauthorized, apperror.ErrInvalidCredentials)
		return
	}

	apperror.LogInfo("User %s logged in successfully", loginReq.Username)
	apperror.WriteSuccess(w, "Login successful", tokenPair)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var refreshReq auth.RefreshTokenRequest
	if err := apperror.ValidateJSON(r, &refreshReq); err != nil {
		apperror.LogWarning("Invalid JSON in refresh token request: %v", err)
		apperror.WriteError(w, http.StatusBadRequest, apperror.ErrInvalidJSON)
		return
	}

	if refreshReq.RefreshToken == "" {
		apperror.WriteError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(refreshReq.RefreshToken)
	if err != nil {
		apperror.LogWarning("Refresh token failed: %v", err)
		apperror.WriteError(w, http.StatusUnauthorized, apperror.ErrInvalidToken)
		return
	}

	apperror.LogInfo("Token refreshed successfully")
	apperror.WriteSuccess(w, "Token refreshed successfully", tokenPair)
}
