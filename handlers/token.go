package handlers

import (
	"jwt_clean/internal/service"
	"jwt_clean/error"
	"net/http"
)

type TokenHandler struct {
	authService *service.Service
}

func NewTokenHandler(authService *service.Service) *TokenHandler {
	return &TokenHandler{authService: authService}
}

func (h *TokenHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.URL.Query().Get("username")
	if username == "" {
		apperror.WriteError(w, http.StatusBadRequest, apperror.ErrMissingRequiredField)
		return
	}

	token, err := h.authService.GenerateToken(username)
	if err != nil {
		apperror.LogWarning("Token generation failed for user %s: %v", username, err)
		apperror.WriteError(w, http.StatusInternalServerError, apperror.ErrTokenGenerationFailed)
		return
	}

	apperror.WriteSuccess(w, "Token generated successfully", map[string]string{"token": token})
}

func (h *TokenHandler) ParseToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		apperror.WriteError(w, http.StatusBadRequest, apperror.ErrMissingToken)
		return
	}

	claims, err := h.authService.ParseToken(tokenString)
	if err != nil {
		apperror.LogWarning("Token parse failed for token: %s, error: %v", tokenString, err)
		apperror.WriteError(w, http.StatusUnauthorized, apperror.ErrInvalidToken)
		return
	}

	response := map[string]any{
		"valid":      true,
		"username":   claims.Subject,
		"issuer":     claims.Issuer,
		"issued_at":  claims.IssuedAt.Time,
		"expires_at": claims.ExpiresAt.Time,
	}

	apperror.WriteSuccess(w, "Token parsed successfully", response)
}
