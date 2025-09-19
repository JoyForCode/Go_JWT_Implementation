package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"jwt_clean/handlers"
	"jwt_clean/internal/middleware"
	"jwt_clean/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	r := mux.NewRouter()

	jwtSecret := []byte(os.Getenv("JWT_SECRET_SIGNING_KEY"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET_SIGNING_KEY must be set in the .env file")
	}

	authService := service.NewAuthService(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	tokenHandler := handlers.NewTokenHandler(authService)
	protectedHandler := handlers.NewProtectedHandler()
	authHandler := handlers.NewAuthHandler(authService)
	//define routes here

	//public routes for checks
	r.HandleFunc("/", testServer).Methods("GET")
	r.HandleFunc("/check", healthCheck).Methods("GET")

	//New Authentication Routes
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/refresh_token", authHandler.RefreshToken).Methods("POST")

	//Legacy token routes for backward compability and tests
	r.HandleFunc("/generate_token", tokenHandler.GenerateToken).Methods("GET")
	r.HandleFunc("/parse_token", tokenHandler.ParseToken).Methods("GET")

	//protected routes
	r.HandleFunc("/dashboard", authMiddleware.RequireAuth(protectedHandler.Dashboard)).Methods("GET")
	r.HandleFunc("/profile", authMiddleware.RequireAuth(protectedHandler.Profile)).Methods("GET")
	r.HandleFunc("/settings", authMiddleware.RequireAuth(protectedHandler.Settings)).Methods("GET")

	server_port := os.Getenv("BACKEND_SERVER_PORT")
	if server_port == "" {
		server_port = ":8080"
		log.Println("BACKEND_SERVER_PORT not set, using default :8080")
	}

	if server_port[0] != ':' {
		server_port = ":" + server_port
	}

	fmt.Printf("Server starting at port %s\n", server_port)

	log.Printf("Starting server on %s", server_port)
	if err := http.ListenAndServe(server_port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func testServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"The server is up and running"}`)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","port":"%s"}`, os.Getenv("BACKEND_SERVER_PORT"))
}
