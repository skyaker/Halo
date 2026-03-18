package middleware

import (
	render "auth_service/internal/render"
	service "auth_service/internal/service"
	"context"
	"net/http"
)

type contextKey string

const UserIdKey contextKey = "user_id"

func AuthMiddleware(s service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session_cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user_session_token := session_cookie.Value

			userId, err := s.CheckToken(r.Context(), user_session_token)
			if err != nil {
				render.HandleError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIdKey, userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
