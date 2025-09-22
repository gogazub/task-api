package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
}

func GenerateUUID() string {
	return uuid.New().String()
}
