package middleware

import (
	"context"
	"jwt_clean/internal/service"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func SetUserInContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UserContextKey, username)
}

func GetUserFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UserContextKey).(string)
	return username, ok
}

type AuthMiddleware struct {
	authService *service.Service
}

func NewAuthMiddleware(authService *service.Service) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, `{"error":"Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.authService.ParseToken(tokenString)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = SetUserInContext(ctx, claims.Subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
