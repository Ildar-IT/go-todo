package authHandler

import (
	"log/slog"
	"net/http"
	_ "todo/docs"
	"todo/internal/entity"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/service"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	log       *slog.Logger
	services  *service.Service
	validator *validator.Validate
}

func NewAuthHandler(log *slog.Logger, services *service.Service, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{log: log, services: services, validator: validator}
}

// Login выполняет вход пользователя
// @Summary Вход пользователя
// @Description Выполняет вход пользователя и возвращает токены доступа и обновления
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials body entity.UserLoginReq true "Данные для входа"
// @Success 200 {object} entity.TokensRes
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 401 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /auth/login [post]
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
		err := h.validator.Struct(user)
		if err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, err.Error(), log)
			return
		}

		resp, err, status := h.services.Auth.Login(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)
	}
}

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя и возвращает токены доступа и обновления
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user body entity.UserRegisterReq true "Данные для регистрации"
// @Success 200 {object} entity.TokensRes
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /auth/register [post]
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

		err := h.validator.Struct(user)
		if err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, err.Error(), log)
			return
		}

		resp, err, status := h.services.Auth.Register(&user)

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		handlers.SendJSONResponse(w, http.StatusOK, resp, log)

	}
}

// UpdateAccessToken обновляет токен доступа
// @Summary Обновление токена доступа
// @Description Обновляет токен доступа с использованием токена обновления
// @Tags auth
// @Security RefreshTokenAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.TokenAccessRes
// @Failure 401 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Security ApiKeyAuth
// @Router /auth/access [post]
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
