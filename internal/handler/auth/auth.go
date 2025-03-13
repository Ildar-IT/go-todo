package authHandler

import (
	"log/slog"
	"net/http"
	"todo/internal/entity"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/service"
)

type AuthHandler struct {
	log      *slog.Logger
	services *service.Service
}

func NewAuthHandler(log *slog.Logger, services *service.Service) *AuthHandler {
	return &AuthHandler{log: log, services: services}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.login"
		log := h.log.With(
			slog.String("op", op),
		)
		var user entity.UserLoginReq
		if err := handlers.DecodeJSONRequest(w, r, &user, log); err != nil {
			return
		}

		resp, err, status := h.services.Auth.Login(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create User:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.register"
		log := h.log.With(
			slog.String("op", op),
		)
		var user entity.UserRegisterReq
		if err := handlers.DecodeJSONRequest(w, r, &user, log); err != nil {
			return
		}
		resp, err, status := h.services.Auth.Register(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Create user:", "tokens", resp)
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)

	}
}

func (h *AuthHandler) UpdateAccessToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.updateTokens"
		log := h.log.With(
			slog.String("op", op),
		)

		claims := r.Context().Value("claims").(*jwtUtils.RefreshClaims)
		token, err, status := h.services.Auth.GenerateAccessToken(claims.UserId, claims.Role)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Update user access token", "token", token)
		handlers.SendJSONResponse(w, http.StatusOK, entity.TokenAccessRes{Access: token}, log)
	}
}
