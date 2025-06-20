package note_handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	models "note_service/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func checkNoteExistence(db *sql.DB, noteId uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM notes WHERE id = $1)`
	err := db.QueryRow(query, noteId).Scan(&exists)
	return exists, err
}

func AddNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noteInfo models.NoteInfo

		if err := json.NewDecoder(r.Body).Decode(&noteInfo); err != nil {
			log.Error().Err(err).Msg("note info json decode")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO notes (id, user_id, type_id, content)
							VALUES ($1, $2, $3, $4)`
		noteId, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg("new note id generation")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(query, noteId, noteInfo.User_id, noteInfo.Type_id, noteInfo.Content)
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

		if err := json.NewDecoder(r.Body).Decode(&noteInfo); err != nil {
			fmt.Print(r.Body)
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

		query := `DELETE FROM notes WHERE id = $1`
		_, err = db.Exec(query, noteInfo.Note_id)
		if err != nil {
			log.Error().Err(err).Msg("note deleting")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Info().Msg("Note deleted successfully")

		w.WriteHeader(http.StatusOK)
	}
}
