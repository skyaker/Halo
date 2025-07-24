package localstore

import (
	"database/sql"
	"fmt"
	"halo/logger"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func getDbPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home db dir: %w", err)
	}
	return filepath.Join(dir, "halo", "db.sqlite"), nil
}

func GetLocalDbConnection() {
	path, err := getDbPath()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("local db path")
	}

	db, err = sql.Open("sqlite", path)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("open local db")
	}

	schema := `
	CREATE TABLE IF NOT EXISTS notes (
  	id UUID PRIMARY KEY,
  	type_id UUID,
  	content TEXT NOT NULL,
  	created_at BIGINT,
  	updated_at BIGINT,
  	ended_at BIGINT,
  	completed BOOLEAN DEFAULT FALSE
	);`
	if _, err := db.Exec(schema); err != nil {
		logger.Logger.Error().Err(err).Msg("create table")
	}

	logger.Logger.Info().Msg("local db connection established")
}
