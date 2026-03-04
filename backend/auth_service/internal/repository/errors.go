package repository

import (
	models "auth_service/internal/models"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return models.ErrAlreadyExists
		case "23503":
			return models.ErrNotFound
		}
	}

	return err
}
