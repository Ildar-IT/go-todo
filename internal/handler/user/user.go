package userHandler

import (
	"log/slog"
	"net/http"
	"todo/internal/entity"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
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
		err := handlers.DecodeJSONRequest(r, &user, log)

		resp, err, status := h.services.User.Register(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create user:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)

	}
}

func (h *UserHandler) UpdateTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.updateTokens"
		log := h.log.With(
			slog.String("op", op),
		)

		claims := r.Context().Value("claims").(*jwtUtils.RefreshClaims)

		if claims.UserId == 0 || claims.Role == "" {
			handlers.SendJSONResponse(w, http.StatusBadRequest, handlers.HTTPErrorRes{Message: "Cannot get payload"}, log)
		}
		resp, err, status := h.services.User.GenerateTokens(claims.UserId, claims.Role)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create user:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)
	}
}
