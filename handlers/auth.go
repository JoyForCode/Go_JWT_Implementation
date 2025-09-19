package handlers

import (
	"encoding/json"
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
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		//http.Error(w, `{"error","Invalid request payload"}`, http.StatusBadRequest)
		apperror.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tokenPair, err := h.authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		//http.Error(w, `{"error":"Invalid Credentials"}`, http.StatusUnauthorized)
		apperror.WriteError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	json.NewEncoder(w).Encode(tokenPair)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var refreshReq auth.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		//http.Error(w, `{"error":"Invalid request payload"}`, http.StatusUnauthorized)
		apperror.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(refreshReq.RefreshToken)
	if err != nil {
		//http.Error(w, `{"error":"Invalid refresh token"}`, http.StatusUnauthorized)
		apperror.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	json.NewEncoder(w).Encode(tokenPair)
}
