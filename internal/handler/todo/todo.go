package todoHandler

import (
	"log/slog"
	"net/http"
	"strconv"
	"todo/internal/entity"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/service"

	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
	log       *slog.Logger
	services  *service.Service
	validator *validator.Validate
}

func NewTodoHandler(log *slog.Logger, services *service.Service, validator *validator.Validate) *TodoHandler {
	return &TodoHandler{log: log, services: services, validator: validator}
}

// @Summary Создать задачу
// @Description Создать новую задачу для текущего пользователя
// @Tags todo
// @Security AccessTokenAuth
// @Accept  json
// @Produce  json
// @Param   todo body entity.TodoCreateReq true "Данные для создания задачи"
// @Success 200 {object} entity.TodoCreateRes
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /todo [post]
func (h *TodoHandler) CreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.todo.CreateTodo"
		log := h.log.With(
			slog.String("op", op),
		)

		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)

		var todo entity.TodoCreateReq

		if err := handlers.DecodeJSONRequest(w, r, &todo, log); err != nil {
			return
		}

		err := h.validator.Struct(todo)
		if err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, err.Error(), log)
			return
		}

		//entity.Todo{Title: todo.Title, Status: todo.Status}
		id, err, status := h.services.Todo.CreateTodo(&entity.Todo{
			User_id:     claims.UserId,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		})
		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}

		log.Info("Create todo id:", id)
		handlers.SendJSONResponse(w, status, entity.TodoCreateRes{Id: id}, log)
	}
}

// @Summary Получить задачу
// @Description Получить задачу для текущего пользователя по id задачи
// @Tags todo
// @Security AccessTokenAuth
// @Accept  json
// @Produce  json
// @Param   todo body entity.TodoCreateReq true "Данные для создания задачи"
// @Success 200 {object} entity.TodoCreateRes
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /todo/{id} [get]
func (h *TodoHandler) GetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get todo id by params
		const op = "handler.todo.GetTodo"
		log := h.log.With(
			slog.String("op", op),
		)
		id, err := strconv.Atoi(r.PathValue("id"))
		if id == 0 || err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, handlers.HTTPErrorRes{Message: "Invalid id"}, log)
			return
		}
		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)
		todo, err, status := h.services.Todo.GetTodo(id, claims.UserId)
		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("GetTodo id:", todo.Id)
		handlers.SendJSONResponse(w, status, todo, log)
	}
}

// UpdateTodo обновляет задачу
// @Summary Обновить задачу
// @Description Обновляет задачу по ID для текущего пользователя
// @Tags todo
// @Security AccessTokenAuth
// @Accept  json
// @Produce  json
// @Param   todo body entity.TodoUpdateReq true "Данные для обновления задачи"
// @Success 200 {object} entity.TodoUpdateRes
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 404 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /todo [patch]
func (h *TodoHandler) UpdateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.todo.UpdateTodo"
		log := h.log.With(
			slog.String("op", op),
		)
		var todoBody entity.TodoUpdateReq

		if err := handlers.DecodeJSONRequest(w, r, &todoBody, log); err != nil {
			return
		}
		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)
		todoBody.UserId = claims.UserId
		todo, err, status := h.services.Todo.UpdateTodo(&todoBody)
		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("UpdateTodo id:", todo.Id)
		handlers.SendJSONResponse(w, status, todo, log)
	}
}

// GetTodos возвращает список задач для текущего пользователя
// @Summary Получить список задач
// @Description Возвращает список задач для текущего пользователя
// @Tags todo
// @Security AccessTokenAuth
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.TodoGetRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /todos [get]
func (h *TodoHandler) GetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get todo id by params
		const op = "handler.todo.GetTodos"
		log := h.log.With(
			slog.String("op", op),
		)
		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)
		todo, err, status := h.services.Todo.GetTodos(claims.UserId)
		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		handlers.SendJSONResponse(w, status, todo, log)
	}
}

// DeleteTodo удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по ID для текущего пользователя
// @Tags todo
// @Security AccessTokenAuth
// @Accept  json
// @Produce  json
// @Param   id path int true "ID задачи"
// @Success 200 {object} int
// @Failure 400 {object} handlers.HTTPErrorRes
// @Failure 404 {object} handlers.HTTPErrorRes
// @Failure 500 {object} handlers.HTTPErrorRes
// @Router /todo/{id} [delete]
func (h *TodoHandler) DeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.todo.DeleteTodo"
		log := h.log.With(
			slog.String("op", op),
		)
		var todoBody entity.TodoUpdateReq

		id, err := strconv.Atoi(r.PathValue("id"))
		if id == 0 || err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, handlers.HTTPErrorRes{Message: "Invalid id"}, log)
		}
		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)
		todoBody.UserId = claims.UserId
		err, status := h.services.Todo.DeleteTodo(id, claims.UserId)
		if err != nil {
			handlers.SendJSONResponse(w, status, handlers.HTTPErrorRes{Message: err.Error()}, log)
			return
		}
		log.Info("Delete todo id:", id)
		handlers.SendJSONResponse(w, status, id, log)
	}
}
