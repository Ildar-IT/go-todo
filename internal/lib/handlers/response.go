package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HTTPErrorRes struct {
	Message string `json:"message"`
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data any, log *slog.Logger) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("Failed to encode response", "error", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
