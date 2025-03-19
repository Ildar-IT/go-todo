package userHandler

import (
	"log/slog"
	"net/http"
	_ "todo/docs"
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

// @Summary Получение всех пользователей
// @Description Получение всех пользователей только для админ роли
// @Tags user
// @Accept  json
// @Produce  json
// @Security AccessTokenAuth
// @Success 200 {object} []entity.User
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 401 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /auth/login [post]
func (h *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.login"
		log := h.log.With(
			slog.String("op", op),
		)

		users, err, status := h.services.User.GetAllUsers()

		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		handlers.SendJSONResponse(w, http.StatusOK, users, log)
	}
}
