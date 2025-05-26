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

func checkUserExistence(db *sql.DB, username *string, email *string) error {
	query := `SELECT user_id 
						FROM users 
						WHERE username = $1 OR email = $2`
	row := db.QueryRow(query, username, email)

	err := row.Scan(new(int))

	return err
}

func checkAdminCredentials(db *sql.DB, username *string, password *string) error {
	query := `SELECT id 
						FROM admins 
						WHERE username = $1 and password = $2`
	row := db.QueryRow(query, username, password)

	err := row.Scan(new(int))

	return err
}

func authInfoCheck(db *sql.DB, email *string, password *string) error {
	var hashedPassword string
	err := checkUserExistence(db, nil, email)
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

func addUserToDB(db *sql.DB, user *models.UserRegisterInfo) (uint, error) {
	var existence bool
	err := checkUserExistence(db, &user.Username, &user.Email)
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

	query := `INSERT INTO users (username, name, email, password_hash, created_at, avatar_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`
	var userId int
	err = db.QueryRow(query, user.Username, user.Name, user.Email, hashedPassword, time.Now().Unix(), "default-avatar.svg").
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

		userId, err := addUserToDB(db, &userData)
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

func AdminLogin(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var adminData models.AdminLoginStruct
		var adminId uint

		if err := json.NewDecoder(r.Body).Decode(&adminData); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := checkAdminCredentials(db, &adminData.Username, &adminData.Password)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		query := `SELECT id 
							FROM admins 
							WHERE username = $1`

		if err = db.QueryRow(query, adminData.Username).Scan(&adminId); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		token, err := CreateAdminToken(redisDb, &adminId)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "admin_token",
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
