package main

import (
	handlers "category_service/internal/handlers"
	dbconn "category_service/internal/repository"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	r := chi.NewRouter()

	r.Post("/api/category", handlers.AddCategory(db))
	r.Delete("/api/category", handlers.DeleteCategory(db))
	r.Get("/api/category", handlers.GetCategory(db))

	log.Info().Msg("category service is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "category service").
			Msg("Server start failed")
	}
}
