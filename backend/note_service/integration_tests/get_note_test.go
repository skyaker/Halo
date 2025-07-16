package note_integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNotes_Success(t *testing.T) {
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

	req, err := http.NewRequest("GET", "http://localhost:8080/api/note?page=1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var notes []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&notes)
	require.NoError(t, err)
	assert.IsType(t, []map[string]interface{}{}, notes)
}

func TestGetNotes_MissingToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/note?page=1", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetNotes_EmptyList(t *testing.T) {
	userCode := uuid.NewString()[:8]
	user := map[string]string{
		"login":    "user_" + userCode,
		"password": "123456",
		"email":    "user" + userCode + "@example.com",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var token string
	for _, c := range resp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	require.NotEmpty(t, token)

	req, err := http.NewRequest("GET", "http://localhost:8080/api/note?page=1", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	client := &http.Client{}
	getResp, err := client.Do(req)
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var notes []map[string]interface{}
	err = json.NewDecoder(getResp.Body).Decode(&notes)
	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestGetNotes_InvalidPage(t *testing.T) {
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

	client := &http.Client{}

	invalidPages := []string{"abc", "-3", ""}

	for _, p := range invalidPages {
		url := "http://localhost:8080/api/note?page=" + p

		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: token,
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var notes []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&notes)
		require.NoError(t, err)

		assert.IsType(t, []map[string]interface{}{}, notes)
	}
}
