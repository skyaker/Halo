package user_handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	models "user_service/internal/models"

	"github.com/rs/zerolog/log"
)

func AddUser(db *sql.DB, data []byte) error {
	var event models.UserRegisterInfo
	if err := json.Unmarshal(data, &event); err != nil {
		log.Error().Err(err).Msg("user-created parse failed")
		return err
	}
	log.Info().Msg("processing user-created")

	query := `INSERT INTO users (id, username, email)
						VALUES ($1, $2, $3)`
	_, err := db.Exec(query, event.User_id, event.Username, event.Email)
	if err != nil {
		log.Error().Err(err).Msg("user info insert error")
		return err
	}

	return nil
}

func DeleteUser(db *sql.DB, data []byte) error {
	var event models.UserDeleteInfo

	if err := json.Unmarshal(data, &event); err != nil {
		log.Error().Err(err).Msg("user-deleted parse failed")
		return err
	}
	log.Info().Msg("processing user-created")

	query := `DELETE FROM users
						WHERE id = $1`
	_, err := db.Exec(query, event.User_id)
	if err != nil {
		log.Error().Err(err).Msg("user info remove error")
		return err
	}

	return nil
}

func CheckUserExistence(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("asdf")
	}
}
