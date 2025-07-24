package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"halo/logger"
	"halo/models"
	"net/http"
)

func SendNoteToService(sessionToken string, noteInfo models.NoteStruct) error {
	noteBody, err := json.Marshal(noteInfo)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("marshal note info")
		return fmt.Errorf("marshal note info: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/note",
		bytes.NewBuffer(noteBody),
	)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("new request")
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("do request")
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error().Msg("request status code")
		return fmt.Errorf("request status code: %d", resp.StatusCode)
	}

	return nil
}
