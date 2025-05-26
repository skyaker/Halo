package main

import (
	"fmt"
	"log"
	"net/http"
	handlers "user_service/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	r := chi.NewRouter()

	r.Post("/users", handlers.Testfunc())

	log.Println("User service is running")
	http.ListenAndServe(":8080", r)
}
