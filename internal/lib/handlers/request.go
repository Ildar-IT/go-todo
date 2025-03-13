package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func DecodeJSONRequest(w http.ResponseWriter, r *http.Request, v any, log *slog.Logger) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Error("Failed to decode request", "error", err.Error())
		SendJSONResponse(w, http.StatusBadRequest, HTTPErrorRes{Message: err.Error()}, log)
		return err
	}
	return nil
}
