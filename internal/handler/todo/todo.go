package todoHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"todo/internal/entity"
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
		var todo entity.TodoCreateReq
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//entity.Todo{Title: todo.Title, Status: todo.Status}
		id, err := h.services.Todo.CreateTodo(&entity.Todo{
			User_id:     todo.User_id,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		})
		if err != nil {
			h.log.Error("Create Todo error", err.Error())
			return
		}

		h.log.Info("Create todo id:", id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(id)
		//json.NewEncoder(w).Encode()
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
