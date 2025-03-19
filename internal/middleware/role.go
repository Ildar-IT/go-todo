package middleware

import (
	"log/slog"
	"net/http"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
)

const (
	AdminRole = "admin"
)

func RoleCheckMiddleware(next http.HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.RoleCheckMiddleware"
		log := logger.With(
			slog.String("op", op),
		)

		claims := r.Context().Value("claims").(*jwtUtils.AccessClaims)

		if claims.Role != AdminRole {
			log.Error("Role check middleware", "error", "Not super user")
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "Not valid token"}, log)
			return
		}

		next.ServeHTTP(w, r)
	}
}
