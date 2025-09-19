package apperror

import (
	"log"
	"time"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	TimeStamp string `json:"timestamp"`
	Code string `json:"code,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
	TimeStamp string `json:"timestamp"`
}

const (
	//Authentication errors
	ErrInvalidCredentials = "Invalid Credentials"
	ErrMissingCredentials = "Username and password are required"
	ErrTokenGenerationFailed = "Failed to generate token"
	
	// Token validation errors
	ErrMissingToken = "Authorization token is required"
	ErrInvalidToken = "Invalid token"
	ErrExpiredToken = "Token has expired"
	ErrMalformedToken = "Token format is invalid"
	ErrTokenNotActive = "Token is not yet active"
	
	// Request errors
	ErrInvalidJSON = "Invalid JSON format"
	ErrMissingRequiredField = "Required field is missing"
	
	// Server errors
	ErrInternalServer = "Internal server error"
)

func getCurrentTimeStamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response:=ErrorResponse{
		Error: message,
		TimeStamp: getCurrentTimeStamp(),
	}

	if err:=json.NewEncoder(w).Encode(response); err!=nil{
		log.Printf("Failed to encode error response: %v", err)
	}

	log.Printf("[ERROR] %d - %s", statusCode, message)
}

func WriteErrorWithCode(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response:=ErrorResponse{
		Error: message,
		Code: code,
		TimeStamp: getCurrentTimeStamp(),
	}

	if err:=json.NewEncoder(w).Encode(response); err!=nil{
		log.Printf("Failed to encode error response: %v", err)
	}

	log.Printf("[ERROR] %d - %s (Code: %s)", statusCode, message, code)
}

func WriteSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response:=SuccessResponse{
		Message: message,
		Data: data,
		TimeStamp: getCurrentTimeStamp(),
	}

	if err:=json.NewEncoder(w).Encode(response); err!=nil{
		log.Printf("Failed to encode success response: %v", err)
	}
}

func ValidateJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func LogInfo(message string, args ...any) {
	log.Printf("[INFO] "+message, args...)
}

func LogWarning(message string, args ...any) {
	log.Printf("[WARNING] "+message, args...)
}
