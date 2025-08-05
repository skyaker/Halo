package category_handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/publicsuffix"

	models "category_service/internal/models"
)

const (
	pageLimit = 10
)

type httpError struct {
	Code  int    `json:"code"`
	Error error  `json:"error"`
	Msg   string `json:"msg"`
}

func checkCategoryExistence(db *sql.DB, userId uuid.UUID, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE user_id = $1 AND name = $2)`
	err := db.QueryRow(query, userId, name).Scan(&exists)
	return exists, err
}

func getUserIdFromToken(r *http.Request) (models.UserInfo, httpError) {
	var userInfo models.UserInfo
	var httpErr httpError

	session_cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			httpErr.Code = http.StatusUnauthorized
			httpErr.Msg = "Unauthorized"
		} else {
			httpErr.Code = http.StatusBadRequest
			httpErr.Msg = "Bad request"
		}
		log.Error().Err(err).Msg("cookie parse error")
		httpErr.Error = err
		return userInfo, httpErr
	}

	user_session_token := session_cookie.Value

	if user_session_token == "" {
		log.Error().Err(fmt.Errorf("session cookie is emtpy"))
		httpErr.Code = http.StatusBadRequest
		httpErr.Error = fmt.Errorf("session cookie is emtpy")
		httpErr.Msg = "Bad request"
		return userInfo, httpErr
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Error().Err(err).Msg("cookie jar init")
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "Internal server error"
		return userInfo, httpErr
	}

	client := http.Client{Jar: jar}

	authUrl := "http://auth_service:8080/api/auth/check_token"

	u, _ := url.Parse(authUrl)
	jar.SetCookies(u, []*http.Cookie{
		{
			Name:  "session_token",
			Value: session_cookie.Value,
		},
	})

	resp, err := client.Get(authUrl)
	if err != nil {
		log.Error().Msg("check token in id parse from token")
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "Internal server error"
		return userInfo, httpErr
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Msg("check token in id parse from token")
		httpErr.Code = resp.StatusCode
		httpErr.Error = fmt.Errorf("check token error")
		return userInfo, httpErr
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msg("auth service check token")
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = fmt.Errorf("auth service check token")
		httpErr.Msg = "Internal server error"
		return userInfo, httpErr
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error().Err(err).Msg("user info json decode")
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "Internal server error"
		return userInfo, httpErr
	}

	return userInfo, httpErr
}

func AddCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, errInfo := getUserIdFromToken(r)
		if errInfo.Error != nil {
			http.Error(w, errInfo.Msg, errInfo.Code)
			return
		}

		var categoryInfo models.CategoryInfo

		if err := json.NewDecoder(r.Body).Decode(&categoryInfo); err != nil {
			log.Error().Err(err).Msg("category info json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if categoryInfo.Name == "" {
			log.Error().Err(fmt.Errorf("category name is empty"))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		exists, err := checkCategoryExistence(db, userInfo.User_id, categoryInfo.Name)
		if err != nil {
			log.Error().Err(err).Msg("check category existence")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if exists {
			log.Error().Err(fmt.Errorf("category name already exists"))
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}

		query := `INSERT INTO categories (id, user_id, name, created_at)
							VALUES ($1, $2, $3, $4)`
		categoryId, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg("new category id generation")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(query, categoryId, userInfo.User_id, categoryInfo.Name, time.Now().Unix())
		if err != nil {
			log.Error().Err(err).Msg("category creating")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, errInfo := getUserIdFromToken(r)
		if errInfo.Error != nil {
			http.Error(w, errInfo.Msg, errInfo.Code)
			return
		}

		var categoryInfo models.CategoryInfo

		if err := json.NewDecoder(r.Body).Decode(&categoryInfo); err != nil {
			log.Error().Err(err).Msg("category id json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		exists, err := checkCategoryExistence(db, userInfo.User_id, categoryInfo.Name)
		if err != nil {
			log.Error().Err(err).Msg("check category existence")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			log.Error().Err(fmt.Errorf("category not found"))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM categories WHERE id = $1 and user_id = $2`
		res, err := db.Exec(query, categoryInfo.Id, userInfo.User_id)
		if err != nil {
			log.Error().Err(err).Msg("category deleting")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Info().Msg("Category deleted successfully")

		w.WriteHeader(http.StatusOK)
	}
}

func GetCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, errInfo := getUserIdFromToken(r)
		if errInfo.Error != nil {
			http.Error(w, errInfo.Msg, errInfo.Code)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page := 1

		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			log.Error().Err(err).Msg("page parse")
		} else {
			page = p
		}

		offset := (page - 1) * pageLimit

		query := `SELECT id, user_id, name, created_at, updated_at
							FROM categories
							WHERE user_id = $1
							ORDER BY created_at DESC
							LIMIT $2 OFFSET $3`
		rows, err := db.Query(query, userInfo.User_id, pageLimit, offset)
		if err != nil {
			log.Error().Err(err).Msg("category receiving")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var categories []models.CategoryInfo

		for rows.Next() {
			var categoryInfo models.CategoryInfo
			var createdAt, updatedAt sql.NullInt64

			err := rows.Scan(
				&categoryInfo.Id,
				&userInfo.User_id,
				&categoryInfo.Name,
				&createdAt,
				&updatedAt,
			)
			if err != nil {
				log.Error().Err(err).Msg("category info scan")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if createdAt.Valid && createdAt.Int64 != 0 {
				categoryInfo.Created_at = createdAt.Int64
			}
			if updatedAt.Valid && updatedAt.Int64 != 0 {
				categoryInfo.Updated_at = updatedAt.Int64
			}
			categories = append(categories, categoryInfo)
		}

		if err := rows.Err(); err != nil {
			log.Error().Err(err).Msg("category receiving")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(categories); err != nil {
			log.Error().Err(err).Msg("failed to write json response")
		}
		return
	}
}
