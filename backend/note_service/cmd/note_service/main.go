package main

import (
	"database/sql"
	"net/http"
	handlers "note_service/internal/handlers"
	dbconn "note_service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	r := chi.NewRouter()

	r.Post("/api/note", handlers.AddNote(db))
	r.Delete("/api/note", handlers.DeleteNote(db))
	r.Get("/api/note", handlers.GetNote(db))

	log.Info().Msg("Note service is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "note service").
			Msg("Server start failed")
	}
}
