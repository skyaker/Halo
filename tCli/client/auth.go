package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"halo/logger"
	"net/http"

	models "halo/models"
)

func Login(login, password string) (string, error) {
	data := models.LoginRequest{
		Login:    login,
		Password: password,
	}

	body, _ := json.Marshal(data)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Client login request")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error().
			Err(fmt.Errorf("status %d", resp.StatusCode)).
			Msg("Client login status code")
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	var session_token string
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			session_token = cookie.Value
		}
	}

	if session_token == "" {
		logger.Logger.Error().Err(err).Msg("Client login cookie")
		return "", err
	}

	return session_token, nil
}
