package config

import (
	"halo/logger"
	"os"
	"path/filepath"
)

func tokenPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("User home dir")
		return "", err
	}
	return filepath.Join(dir, "halo", "token"), nil
}

func SaveToken(token string) error {
	path, err := tokenPath()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Token path")
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Mkdir token path")
		return err
	}

	return os.WriteFile(path, []byte(token), 0600)
}

func LoadToken() (string, error) {
	path, err := tokenPath()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Token path")
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Read token file")
		return "", err
	}
	return string(data), nil
}
