package main

import (
	"log"
	"net/http"
	handlers "user_service/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/api/user", handlers.Testfunc())
	// r.Get
	// r.Delete()
	// r.Put

	log.Println("User service is running")
	http.ListenAndServe(":8100", r)
}
