package main

import (
	"database/sql"
	"log"
	"net/http"
	handlers "user_service/internal/handlers"
	dbconn "user_service/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()

	r := chi.NewRouter()

	r.Post("/api/user", handlers.AddUser())
	r.Get("/api/user/existence", handlers.CheckUserExistence(db))
	// r.Get
	// r.Delete()
	// r.Put

	log.Println("User service is running")
	http.ListenAndServe(":8080", r)
}
