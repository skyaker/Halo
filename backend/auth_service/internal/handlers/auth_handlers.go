package auth_handlers

import (
	models "auth_service/internal/models"
	service "auth_service/internal/service"
	"encoding/json"
	"errors"
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
			if errors.Is(err, models.ErrAlreadyExists) {
				http.Error(w, "User already exists", http.StatusBadRequest)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

// func (h *AuthHandler) HandleDelete() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session_cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			log.Error().Err(err).Msg("Getting cookie")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		user_session_token := session_cookie.Value
//
// 		if user_session_token == "" {
// 			log.Error().Msg("Session token is empty")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		user_id, errInfo := extractUserIdFromToken(user_session_token)
//
// 		if errInfo != nil {
// 			if strings.Contains(errInfo.Error(), "unauthorized") {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			} else if strings.Contains(errInfo.Error(), "bad request") {
// 				http.Error(w, "Bad request", http.StatusBadRequest)
// 			} else {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			}
// 			return
// 		}
//
// 		exists, err := checkUserExistenceById(h.db, &user_id)
// 		if err != nil {
// 			log.Error().Err(err).Msg("check user existence")
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if !exists {
// 			log.Error().Err(fmt.Errorf("user not found")).Msg("user doesn't exist")
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
//
// 		query := `DELETE FROM auth_credentials
// 							WHERE user_id = $1`
//
// 		res, err := h.db.Exec(query, user_id)
// 		if err != nil {
// 			log.Error().Err(err).Msg("deleting user credentials")
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
// 			log.Error().Msg("User not found")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		err = h.writer.WriteMessages(r.Context(), kafka.Message{
// 			Topic: "user-deleted",
// 			Key:   []byte(user_id.String()),
// 			Value: []byte(user_id.String()),
// 		})
// 		if err != nil {
// 			log.Error().Err(err).Msg("kafka user-deleted message error")
// 		} else {
// 			log.Info().Msg("kafka message user-deleted sent")
// 		}
//
// 		redisRes, err := h.redisDb.Del(context.Background(), fmt.Sprintf("user:%v", user_id)).
// 			Result()
// 		if err != nil {
// 			log.Error().Err(err).Msg("redis token removal")
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if redisRes == 0 {
// 			log.Warn().Msg("Token not found")
// 		}
//
// 		log.Info().Msg("User credentials deleted successfully")
//
// 		w.WriteHeader(http.StatusOK)
// 	}
// }

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
// func checkUserExistenceById(db *sql.DB, user_id *uuid.UUID) (bool, error) {
// 	var exists bool
// 	query := `SELECT EXISTS (SELECT 1 FROM auth_credentials WHERE user_id = $1)`
// 	err := db.QueryRow(query, user_id.String()).Scan(&exists)
// 	return exists, err
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
