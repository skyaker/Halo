package auth_handlers

import (
	models "auth_service/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
)

func checkUserExistence(db *sql.DB, login *string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM auth_credentials WHERE login = $1)`
	err := db.QueryRow(query, login).Scan(&exists)
	return exists, err
}

func authInfoCheck(db *sql.DB, login *string, password *string) error {
	var hashedPassword string

	exists, err := checkUserExistence(db, login)
	if err != nil {
		log.Error().Err(err).Msg("check user existence")
		return fmt.Errorf("internal server error: %v", err.Error())
	}
	if !exists {
		log.Error().Err(fmt.Errorf("invalid credentials"))
		return fmt.Errorf("unauthorized: invalid credentials")
	}

	query := `SELECT password_hash
						FROM auth_credentials
						WHERE login = $1`

	if err = db.QueryRow(query, login).Scan(&hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("invalid credentials"))
			return fmt.Errorf("unauthorized: invalid credentials")
		} else {
			log.Error().Err(err).Msg("extracting user password from db")
			return fmt.Errorf("internal server error: %v", err.Error())
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(*password))
	if err != nil {
		log.Error().Err(err).Msg("password comparing")
		return fmt.Errorf("unauthorized: invalid credentials")
	}

	log.Info().Msg("credentials check passed")
	return nil
}

func addUserCredentials(db *sql.DB, user *models.UserRegisterInfo) (uuid.UUID, error) {
	existence, err := checkUserExistence(db, &user.Login)
	if err != nil {
		log.Error().Err(err).Msg("check user existence")
		return uuid.UUID{}, err
	}

	if existence {
		return uuid.UUID{}, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error hashing password")
	}

	query := `INSERT INTO auth_credentials (user_id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)`
	userId, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	_, err = db.Exec(query, userId, user.Login, hashedPassword, time.Now().Unix())
	if err != nil {
		return uuid.UUID{}, err
	}

	return userId, nil
}

func RegisterUser(db *sql.DB, redisDb *redis.Client, writer *kafka.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserRegisterInfo

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			log.Error().Err(err).Msg("register json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		userId, err := addUserCredentials(db, &userData)
		if err != nil {
			if err.Error() == "user already exists" {
				log.Error().Err(err).Msg("User already exists")
				http.Error(w, "User already exists", http.StatusBadRequest)
				return
			} else if err.Error() == "database error" {
				log.Error().Err(err).Msg("Add credentials database error")
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			} else {
				log.Error().Err(err).Msg("adding user credentials")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		log.Info().Msg("User credentials added successfully")

		event := map[string]any{
			"id":       userId,
			"username": "",
			"email":    "",
		}

		if userData.Username != "" {
			event["username"] = userData.Username
		}

		if userData.Email != "" {
			event["email"] = userData.Email
		}

		data, _ := json.Marshal(event)

		err = writer.WriteMessages(r.Context(), kafka.Message{
			Key:   []byte(userId.String()),
			Value: data,
		})
		if err != nil {
			log.Error().Err(err).Msg("kafka user-created message error")
		} else {
			log.Info().Msg("kafka message user-created sent")
		}

		token, err := CreateToken(redisDb, &userId)
		if err != nil {
			log.Error().
				Err(err).
				Str("process", "user token creating")
			return
		}
		log.Info().Msg("Session token created successfully")

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

func DeleteUser(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserDeletedEvent

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			log.Error().Err(err).Msg("delete request decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		exists, err := checkUserExistence(db, &userData.Login)
		if err != nil {
			log.Error().Err(err).Msg("check user existence")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			log.Error().Err(fmt.Errorf("user not found"))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		query := `DELETE FROM auth_credentials
							WHERE login = $1`

		_, err = db.Exec(query, userData.Login)
		if err != nil {
			log.Error().Err(err).Msg("deleting user credentials")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Info().Msg("User credentials deleted successfully")

		w.WriteHeader(http.StatusOK)
	}
}

func Login(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserLogin
		var userId uuid.UUID

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			log.Error().Err(err).Msg("login request decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := authInfoCheck(db, &userData.Login, &userData.Password)
		if err != nil {
			if strings.Contains(err.Error(), "unauthorized") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		query := `SELECT user_id
							FROM auth_credentials
							WHERE login = $1`

		if err = db.QueryRow(query, userData.Login).Scan(&userId); err != nil {
			if err == sql.ErrNoRows {
				log.Error().Err(err).Msg("user id was not found by login")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				log.Error().Err(err).Msg("extracting user id by login")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		token, err := CreateToken(redisDb, &userId)
		if err != nil {
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
