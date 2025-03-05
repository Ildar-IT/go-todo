package todoHandler

import (
	"log/slog"
	"net/http"
	servicesTodo "todo/internal/services/todo"
)

type TodoHandler struct {
	service *services.CartService
}

func NewCartHandler(log *slog.Logger, router *http.ServeMux, service *servicesTodo.TodoService) {
	cart := TodoHandler{service: service}
	router.HandleFunc("GET /cart", cart.GetCart())
}

func (s *TodoHandler) GetCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text")
		w.WriteHeader(200)
		w.Write([]byte("На получи 1"))
		w.Write([]byte("На получи 2"))
		w.Write([]byte("На получи 3"))
		//json.NewEncoder(w).Encode()
	}
}
