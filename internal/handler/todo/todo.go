package todoHandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"todo/internal/entity"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/service"
)

type TodoHandler struct {
	log      *slog.Logger
	services *service.Service
}

func NewTodoHandler(log *slog.Logger, services *service.Service) *TodoHandler {
	return &TodoHandler{log: log, services: services}
}

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

func (h *TodoHandler) GetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get todo id by params
		const op = "handler.todo.GetTodo"
		log := h.log.With(
			slog.String("op", op),
		)
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if id == 0 || err != nil {
			handlers.SendJSONResponse(w, http.StatusBadRequest, handlers.HTTPErrorRes{Message: "Invalid id"}, log)
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
func (h *TodoHandler) DeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//delete todo id by params
	}
}
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

func (h *TodoHandler) GetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get todos by user id
	}
}

// func (s *TodoHandler) GetCart() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var todo TodoCreateRequest

// 		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		err := s.service.CreateTodo(todo)

// 		// w.Header().Set("Content-Type", "text")
// 		// w.WriteHeader(200)
// 		// w.Write([]byte("На получи 1"))
// 		// w.Write([]byte("На получи 2"))
// 		// w.Write([]byte("На получи 3"))
// 		//json.NewEncoder(w).Encode()
// 	}
// }
