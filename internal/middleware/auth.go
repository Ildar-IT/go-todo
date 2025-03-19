package middleware

import (
	"context"
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
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid authorization header"}, log)
			return
		}
		if len(authParts[1]) == 0 {
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid authorization header"}, log)
			return
		}
		claims, err := jwt.ValidateAccessToken(authParts[1])

		if err != nil {
			log.Error("Token not valid", "error", err.Error())
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "Not valid token"}, log)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
	}
}
func RefreshTokenMiddleware(next http.HandlerFunc, jwt *jwtUtils.Jwt, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.RefreshTokenMiddleware"
		log := logger.With(
			slog.String("op", op),
		)
		token := r.Header.Get("Authorization")
		authParts := strings.Split(token, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid authorization header"}, log)
			return
		}
		if len(authParts[1]) == 0 {
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "invalid authorization header"}, log)
			return
		}
		claims, err := jwt.ValidateRefreshToken(authParts[1])

		if err != nil {
			log.Error("Token not valid", "error", err.Error())
			handlers.SendJSONResponse(w, http.StatusUnauthorized, handlers.HTTPErrorRes{Message: "Not valid token"}, log)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
	}
}
