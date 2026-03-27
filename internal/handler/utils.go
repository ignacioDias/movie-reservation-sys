package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteResponseWithEncoder(w http.ResponseWriter, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
