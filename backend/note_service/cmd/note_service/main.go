package main

import (
	"database/sql"
	"net/http"
	handlers "note_service/internal/handlers"
	dbconn "note_service/internal/repository"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	_, envSt := os.LookupEnv("NOTE_POSTGRES")
	if !envSt {
		log.Fatal().Msg("Note host name not found")
	}

	_, envSt = os.LookupEnv("NOTE_PG_PORT")
	if !envSt {
		log.Fatal().Msg("Note postgres port not found")
	}

	_, envSt = os.LookupEnv("POSTGRES_USER")
	if !envSt {
		log.Fatal().Msg("Note postgres user not found")
	}

	_, envSt = os.LookupEnv("POSTGRES_PASSWORD")
	if !envSt {
		log.Fatal().Msg("Note postgres password not found")
	}

	_, envSt = os.LookupEnv("NOTE_DB")
	if !envSt {
		log.Fatal().Msg("Note db name not found")
	}

	r := chi.NewRouter()

	r.Post("/api/note", handlers.AddNote(db))

	log.Info().Msg("Note service is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "note service").
			Msg("Server start failed")
	}
}
