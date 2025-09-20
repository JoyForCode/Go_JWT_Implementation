package apperror

import (
	"log"
	"time"
	"encoding/json"
	"net/http"
	"errors"
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

var (
	//Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMissingCredentials = errors.New("username and password are required")
	ErrTokenGenerationFailed = errors.New("failed to generate token")
	
	// Token validation errors
	ErrMissingToken = errors.New("authorization token is required")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrMalformedToken = errors.New("token format is invalid")
	ErrTokenNotActive = errors.New("token is not yet active")
	
	// Request errors
	ErrInvalidJSON = errors.New("invalid json format")
	ErrMissingRequiredField = errors.New("required field is missing")
	
	// Server errors
	ErrInternalServer = errors.New("internal server error")
)

func getCurrentTimeStamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response:=ErrorResponse{
		Error: err.Error(),
		TimeStamp: getCurrentTimeStamp(),
	}

	if err:=json.NewEncoder(w).Encode(response); err!=nil{
		log.Printf("Failed to encode error response: %v", err)
	}

	log.Printf("[ERROR] %d - %s", statusCode, err)
}

func WriteAppError(w http.ResponseWriter, statusCode int, err error, code string) {
	if code==""{
		WriteError(w, statusCode, err)
		return
	}
	WriteErrorWithCode(w, statusCode, err, code)
}

func WriteErrorWithCode(w http.ResponseWriter, statusCode int, err error, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response:=ErrorResponse{
		Error: err.Error(),
		Code: code,
		TimeStamp: getCurrentTimeStamp(),
	}

	if err:=json.NewEncoder(w).Encode(response); err!=nil{
		log.Printf("Failed to encode error response: %v", err)
	}

	log.Printf("[ERROR] %d - %s (Code: %s)", statusCode, err, code)
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
