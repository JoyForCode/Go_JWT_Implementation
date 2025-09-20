package handlers

import (
	"encoding/json"
	"jwt_clean/internal/middleware"
	"net/http"
	"time"
)

type ProtectedHandler struct{}

func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

func (h *ProtectedHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		username = "unknown"
	}

	response := map[string]interface{}{
		"message":     "Welcome to your dashboard!",
		"user":        username,
		"access_time": time.Now(),
		"permissions": []string{"read", "write", "delete"},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ProtectedHandler) Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username, _ := middleware.GetUserFromContext(r.Context())
	response := map[string]any{
		"username": username,
		"email":    username + "@example.com",
		"role":     "user",
		"profile":  "This is a protected profile endpoint",
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ProtectedHandler) Settings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username, _ := middleware.GetUserFromContext(r.Context())
	response := map[string]any{
		"user":    username,
		"message": "User settings - only accessible with vaild JWT",
		"settings": map[string]string{
			"theme":    "dark",
			"language": "en",
			"timezone": "IST",
		},
	}

	json.NewEncoder(w).Encode(response)
}
