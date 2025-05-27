package auth_handlers

import (
	models "auth_service/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func checkUserExistence(db *sql.DB, login *string) error {
	query := `SELECT user_id 
						FROM auth_credentials
						WHERE login = $1`
	row := db.QueryRow(query, login)

	err := row.Scan(new(int))

	return err
}

func authInfoCheck(db *sql.DB, email *string, password *string) error {
	var hashedPassword string
	err := checkUserExistence(db, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("unauthorized: invalid credentials")
		} else {
			return fmt.Errorf("internal server error: %v", err.Error())
		}
	}

	query := `SELECT password_hash
						FROM users
						WHERE email = $1`

	if err = db.QueryRow(query, email).Scan(&hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("unauthorized: invalid credentials")
		} else {
			return fmt.Errorf("internal server error: %v", err.Error())
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(*password))
	if err != nil {
		return fmt.Errorf("unauthorized: invalid credentials")
	}

	return nil
}

func addUserCredentials(db *sql.DB, user *models.UserRegisterInfo) (uint, error) {
	var existence bool
	err := checkUserExistence(db, &user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			existence = false
		} else {
			return 0, fmt.Errorf("database error")
		}
	}

	if existence {
		return 0, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hashing password")
	}

	query := `INSERT INTO auth_credentials (login, password_hash, created_at) VALUES ($1, $2, $3) RETURNING user_id`
	var userId int
	err = db.QueryRow(query, user.Login, hashedPassword, time.Now().Unix()).
		Scan(&userId)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func RegisterUser(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserRegisterInfo

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		userId, err := addUserCredentials(db, &userData)
		if err != nil {
			if err.Error() == "user already exists" {
				http.Error(w, "User already exists", http.StatusBadRequest)
				return
			} else if err.Error() == "database error" {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
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

func Login(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.UserLogin
		var userId uint

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := authInfoCheck(db, &userData.Email, &userData.Password)
		if err != nil {
			if strings.Contains(err.Error(), "unauthorized") {
				// log (?)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		query := `SELECT user_id 
							FROM users 
							WHERE email = $1`

		if err = db.QueryRow(query, userData.Email).Scan(&userId); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
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
