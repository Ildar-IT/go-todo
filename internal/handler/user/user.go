package userHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todo/internal/entity"
	"todo/internal/lib/handlers"
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
		err := handlers.DecodeJSONRequest(r, &user, log)
		if err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, handlers.HTTPErrorRes{Message: err.Error()}, log)
		}

		resp, err, status := h.services.User.Login(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create User:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)
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
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create user:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)

	}
}

// func (h *UserHandler) UpdateTokens() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handler.user.register"
// 		log := h.log.With(
// 			slog.String("op", op),
// 		)
// 	}
// }
