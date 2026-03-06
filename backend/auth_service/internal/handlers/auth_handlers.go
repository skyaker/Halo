package auth_handlers

import (
	models "auth_service/internal/models"
	service "auth_service/internal/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) HandleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserRegisterInfo

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			log.Error().Err(err).Msg("register json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		token, err := h.service.RegisterUser(r.Context(), userData)
		if err != nil {
			log.Error().Err(err).Msg("user registration failed")
			handleError(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Expires:  time.Now().Add(6 * 30 * 24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandler) HandleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Error().Err(err).Msg("Getting cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user_session_token := session_cookie.Value

		err = h.service.DeleteUser(r.Context(), user_session_token)
		if err != nil {
			log.Error().Err(err).Msg("user deletion failed")
			handleError(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandler) CheckToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Error().Err(err).Msg("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user_session_token := session_cookie.Value

		userId, err := h.service.CheckToken(r.Context(), user_session_token)
		if err != nil {
			log.Error().Err(err).Msg("check token failed")
			handleError(w, err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userId,
		})
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserLogin

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			log.Error().Err(err).Msg("login request decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		token, err := h.service.Login(r.Context(), userData)
		if err != nil {
			log.Error().Err(err).Msg("user login failed")
			handleError(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Expires:  time.Now().Add(6 * 30 * 24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
	}
}
