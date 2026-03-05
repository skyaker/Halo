package auth_handlers

import (
	models "auth_service/internal/models"
	"errors"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		http.Error(w, "Not found", http.StatusNotFound)
	case errors.Is(err, models.ErrAlreadyExists):
		http.Error(w, "User already exists", http.StatusConflict)
	case errors.Is(err, models.ErrInternal):
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	case errors.Is(err, models.ErrInvalidRequest):
		http.Error(w, "Invalid request", http.StatusBadRequest)
	case errors.Is(err, models.ErrIncorrectPassword):
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
	case errors.Is(err, models.ErrInvalidToken):
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
