package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func DecodeJSONRequest(r *http.Request, v any, log *slog.Logger) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Error("Failed to decode request", "error", err.Error())
		return err
	}

	return nil
}
