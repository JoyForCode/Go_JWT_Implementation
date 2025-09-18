package handlers

import (
	"encoding/json"
	"jwt_clean/internal/service"
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
		http.Error(w, `{"error":"Username required"}`, http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(username)
	if err != nil {
		http.Error(w, `{"error":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *TokenHandler) ParseToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, `{"error":"Token required"}`, http.StatusBadRequest)
		return
	}

	claims, err := h.authService.ParseToken(tokenString)
	if err != nil {
		http.Error(w, `{"error":"Invalid Token"}`, http.StatusUnauthorized)
		return
	}

	response := map[string]any{
		"valid":      true,
		"username":   claims.Subject,
		"issuer":     claims.Issuer,
		"issued_at":  claims.IssuedAt.Time,
		"expires_at": claims.ExpiresAt.Time,
	}

	json.NewEncoder(w).Encode(response)
}
