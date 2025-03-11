package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"todo/internal/lib/handlers"
	jwtUtils "todo/internal/lib/jwt"
)

func AuthMiddleware(next http.HandlerFunc, jwt *jwtUtils.Jwt, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.authMiddleware"
		log := logger.With(
			slog.String("op", op),
		)
		token := r.Header.Get("Authorization")
		authParts := strings.Split(token, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid auth header"}, log)
			return
		}
		if len(authParts[1]) == 0 {
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid auth header"}, log)
			return
		}
		// claims, err := jwt.ValidateAccessToken(authParts[1])
		// log.Info("USER ID", "id", claims["user_id"].(int))
		// log.Info("USER ID", "id", claims["user_id"].(string))
		// if err != nil {
		// 	log.Error(err.Error())
		// 	handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "Not valid token"}, log)
		// 	return
		// }
		//r = r.WithContext(context.WithValue(r.Context(), "userID"))
		next.ServeHTTP(w, r)
	}
}
