package note_integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddNote_Success(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}
	loginBody, _ := json.Marshal(creds)

	loginResp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(loginBody),
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()

	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var sessionToken string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			sessionToken = c.Value
			break
		}
	}
	require.NotEmpty(t, sessionToken, "session_token not found after login")

	note := map[string]interface{}{
		"category_id": uuid.New().String(),
		"content":     "Test note from login-based test",
	}
	noteBody, _ := json.Marshal(note)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/note", bytes.NewBuffer(noteBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAddNote_InvalidJSON(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}
	loginBody, _ := json.Marshal(creds)

	loginResp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(loginBody),
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var sessionToken string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			sessionToken = c.Value
			break
		}
	}
	require.NotEmpty(t, sessionToken)

	invalidJSON := `{"category_id": "invalid", "content": "broken"`

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/note",
		strings.NewReader(invalidJSON),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAddNote_MissingToken(t *testing.T) {
	note := map[string]interface{}{
		"category_id": uuid.New().String(),
		"content":     "Note without token",
	}
	body, err := json.Marshal(note)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/note", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAddNote_EmptyBody(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}
	loginBody, _ := json.Marshal(creds)

	loginResp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(loginBody),
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var sessionToken string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			sessionToken = c.Value
			break
		}
	}
	require.NotEmpty(t, sessionToken)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/note", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
