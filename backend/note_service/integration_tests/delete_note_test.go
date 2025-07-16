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

func TestDeleteNote_Success(t *testing.T) {
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

	var token string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	require.NotEmpty(t, token)

	typeID := uuid.New()
	note := map[string]interface{}{
		"type_id": typeID.String(),
		"content": "Note to be deleted",
	}
	noteBody, _ := json.Marshal(note)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/note", bytes.NewBuffer(noteBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	getReq, _ := http.NewRequest("GET", "http://localhost:8080/api/note?page=1", nil)
	getReq.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	getResp, err := client.Do(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)

	var notes []map[string]interface{}
	err = json.NewDecoder(getResp.Body).Decode(&notes)
	require.NoError(t, err)
	require.NotEmpty(t, notes)

	noteID := notes[0]["note_id"].(string)

	delBody, _ := json.Marshal(map[string]string{
		"note_id": noteID,
	})

	delReq, _ := http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/note",
		bytes.NewBuffer(delBody),
	)
	delReq.Header.Set("Content-Type", "application/json")
	delReq.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	delResp, err := client.Do(delReq)
	require.NoError(t, err)
	defer delResp.Body.Close()

	assert.Equal(t, http.StatusOK, delResp.StatusCode)
}

func TestDeleteNote_NonexistentID(t *testing.T) {
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

	var token string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	require.NotEmpty(t, token)

	nonexistentID := uuid.New().String()

	delBody, _ := json.Marshal(map[string]string{
		"note_id": nonexistentID,
	})

	req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/note", bytes.NewBuffer(delBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteNote_InvalidUUID(t *testing.T) {
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

	var token string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	require.NotEmpty(t, token)

	body := []byte(`{"note_id": "not-a-uuid"}`)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/note", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteNote_InvalidJSON(t *testing.T) {
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

	var token string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	require.NotEmpty(t, token)

	badJSON := `{"note_id":`

	req, err := http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/note",
		strings.NewReader(badJSON),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteNote_MissingToken(t *testing.T) {
	body, err := json.Marshal(map[string]string{
		"note_id": uuid.New().String(),
	})
	require.NoError(t, err)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/note", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
