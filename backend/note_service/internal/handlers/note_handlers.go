package note_handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	models "note_service/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

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
