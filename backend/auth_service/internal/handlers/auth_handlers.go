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
		if userData.Login == "" || userData.Password == "" {
			log.Error().Msg("empty field")
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

		if user_session_token == "" {
			log.Error().Msg("Session token is empty")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

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

		if user_session_token == "" {
			log.Error().Msg("Session token is empty")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

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

// func (h *AuthHandler) Login() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var userData models.UserLogin
// 		var userId uuid.UUID
//
// 		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
// 			log.Error().Err(err).Msg("login request decode")
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
//
// 		if userData.Login == "" || userData.Password == "" {
// 			log.Error().Err(fmt.Errorf("empty field"))
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
//
// 		err := authInfoCheck(h.db, &userData.Login, &userData.Password)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "unauthorized") {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			} else {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			}
// 			return
// 		}
//
// 		query := `SELECT user_id
// 							FROM auth_credentials
// 							WHERE login = $1`
//
// 		if err = h.db.QueryRow(query, userData.Login).Scan(&userId); err != nil {
// 			if err == sql.ErrNoRows {
// 				log.Error().Err(err).Msg("user id was not found by login")
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			} else {
// 				log.Error().Err(err).Msg("extracting user id by login")
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			}
// 		}
//
// 		token, err := CreateToken(h.redisDb, &userId)
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		http.SetCookie(w, &http.Cookie{
// 			Name:     "session_token",
// 			Value:    token,
// 			HttpOnly: true,
// 			Secure:   false,
// 			SameSite: http.SameSiteLaxMode,
// 			Path:     "/",
// 			Expires:  time.Now().Add(6 * 30 * 24 * time.Hour),
// 		})
//
// 		w.WriteHeader(http.StatusOK)
// 	}
// }
//
// func authInfoCheck(db *sql.DB, login *string, password *string) error {
// 	var hashedPassword string
//
// 	exists, err := checkUserExistenceByLogin(db, login)
// 	if err != nil {
// 		log.Error().Err(err).Msg("check user existence")
// 		return fmt.Errorf("internal server error: %v", err.Error())
// 	}
// 	if !exists {
// 		log.Error().Err(fmt.Errorf("invalid credentials"))
// 		return fmt.Errorf("unauthorized: invalid credentials")
// 	}
//
// 	query := `SELECT password_hash
// 						FROM auth_credentials
// 						WHERE login = $1`
//
// 	if err = db.QueryRow(query, login).Scan(&hashedPassword); err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Error().Err(fmt.Errorf("invalid credentials"))
// 			return fmt.Errorf("unauthorized: invalid credentials")
// 		} else {
// 			log.Error().Err(err).Msg("extracting user password from db")
// 			return fmt.Errorf("internal server error: %v", err.Error())
// 		}
// 	}
//
// 	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(*password))
// 	if err != nil {
// 		log.Error().Err(err).Msg("password comparing")
// 		return fmt.Errorf("unauthorized: invalid credentials")
// 	}
//
// 	log.Info().Msg("credentials check passed")
// 	return nil
// }
