package note_handlers

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

	models "note_service/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/publicsuffix"
)

const (
	pageLimit = 10
)

type httpError struct {
	Code  int    `json:"code"`
	Error error  `json:"error"`
	Msg   string `json:"msg"`
}

func checkNoteExistence(db *sql.DB, noteId uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM notes WHERE id = $1)`
	err := db.QueryRow(query, noteId).Scan(&exists)
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
		log.Error().Err(err).Msg("auth service check token")
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "Internal server error"
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

func AddNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, errInfo := getUserIdFromToken(r)
		if errInfo.Error != nil {
			http.Error(w, errInfo.Msg, errInfo.Code)
			return
		}

		var noteInfo models.NoteInfo

		if err := json.NewDecoder(r.Body).Decode(&noteInfo); err != nil {
			log.Error().Err(err).Msg("note info json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO notes (id, user_id, category_id, content, created_at)
							VALUES ($1, $2, $3, $4, $5)`
		noteId, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg("new note id generation")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(
			query,
			noteId,
			userInfo.User_id,
			noteInfo.Category_id,
			noteInfo.Content,
			time.Now().Unix(),
		)
		if err != nil {
			log.Error().Err(err).Msg("note creating")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Info().Msg("Note inserted successfully")

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noteInfo models.NoteDeleteInfo

		userInfo, errInfo := getUserIdFromToken(r)
		if errInfo.Error != nil {
			http.Error(w, errInfo.Msg, errInfo.Code)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&noteInfo); err != nil {
			log.Error().Err(err).Msg("note id json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		exists, err := checkNoteExistence(db, noteInfo.Note_id)
		if err != nil {
			log.Error().Err(err).Msg("check note existence")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			log.Error().Err(fmt.Errorf("note not found"))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM notes WHERE id = $1 and user_id = $2`
		res, err := db.Exec(query, noteInfo.Note_id, userInfo.User_id)
		if err != nil {
			log.Error().Err(err).Msg("note deleting")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Info().Msg("Note deleted successfully")

		w.WriteHeader(http.StatusOK)
	}
}

func GetNote(db *sql.DB) http.HandlerFunc {
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

		query := `SELECT id, user_id, category_id, content, created_at, updated_at, ended_at, completed
							FROM notes
							WHERE user_id = $1
							ORDER BY created_at DESC
							LIMIT $2 OFFSET $3`

		rows, err := db.Query(query, userInfo.User_id, pageLimit, offset)
		if err != nil {
			log.Error().Err(err).Msg("note receiving")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var notes []models.NoteInfo

		for rows.Next() {
			var noteInfo models.NoteInfo
			var categoryId sql.NullString
			var createdAt, updatedAt, endedAt sql.NullInt64

			err := rows.Scan(
				&noteInfo.Id,
				&userInfo.User_id,
				&categoryId,
				&noteInfo.Content,
				&createdAt,
				&updatedAt,
				&endedAt,
				&noteInfo.Completed,
			)
			if err != nil {
				log.Error().Err(err).Msg("note info scan")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if categoryId.Valid {
				noteInfo.Category_id, err = uuid.Parse(categoryId.String)
				if err != nil {
					log.Error().Err(err).Msg("category id parse")
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
			if createdAt.Valid && createdAt.Int64 != 0 {
				noteInfo.Created_at = createdAt.Int64
			}
			if updatedAt.Valid && updatedAt.Int64 != 0 {
				noteInfo.Updated_at = updatedAt.Int64
			}
			if endedAt.Valid && endedAt.Int64 != 0 {
				noteInfo.Ended_at = endedAt.Int64
			}
			notes = append(notes, noteInfo)
		}

		if err := rows.Err(); err != nil {
			log.Error().Err(err).Msg("note receiving")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(notes); err != nil {
			log.Error().Err(err).Msg("failed to write json response")
		}
		return
		// check author valid
	}
}
