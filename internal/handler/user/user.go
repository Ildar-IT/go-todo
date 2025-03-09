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
		// var user entity.UserCreateReq
		// if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// id, err := h.services.User.Login(&user)

		// if err != nil {
		// 	h.log.Error("Create User error", err.Error())
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// h.log.Info("Create User id:", id)
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(200)
		// json.NewEncoder(w).Encode(id)
	}
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user entity.UserRegisterReq
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.services.User.Register(&user)

		if err != nil {
			h.log.Error("Create User error", err.Error())
			http.Error(w, "Create User error", http.StatusInternalServerError)
			return
		}
		h.log.Info("Create User:", resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(resp)
	}
}
