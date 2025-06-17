package main

import (
	"database/sql"
	"net/http"
	handlers "user_service/internal/handlers"
	user_kafka "user_service/internal/kafka"
	dbconn "user_service/internal/repository"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	go user_kafka.RunKafkaListener(db)

	r := chi.NewRouter()

	r.Get("/api/user/existence", handlers.CheckUserExistence(db))

	log.Info().Msg("User service is running")
	http.ListenAndServe(":8080", r)
}
