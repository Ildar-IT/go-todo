package userHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todo/internal/entity"
	"todo/internal/service"
)

type UserHandler struct {
	log      *slog.Logger
	services *service.Service
}

func NewUserHandler(log *slog.Logger, services *service.Service) *UserHandler {
	return &UserHandler{log: log, services: services}
}

func (h *UserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.login"
		log := h.log.With(
			slog.String("op", op),
		)
		var user entity.UserLoginReq
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error("Failed to Decode response", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err, status := h.services.User.Login(&user)

		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		log.Info("Create User:", "tokens", resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response", "error", err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.register"
		log := h.log.With(
			slog.String("op", op),
		)
		var user entity.UserRegisterReq
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err, status := h.services.User.Register(&user)

		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		log.Info("Create user:", "tokens", resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response", "error", err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
