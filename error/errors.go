package apperror

import (
	"time"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	TimeStamp string `json:"timestamp"`
}

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
	json.NewEncoder(w).Encode(response)
}
