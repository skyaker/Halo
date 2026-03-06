package auth_integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CheckToken_Success(t *testing.T) {
	login := "user_" + uuid.NewString()[:8]
	user := map[string]string{
		"login":    login,
		"password": "password123",
		"email":    login + "@example.com",
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
	require.NotEmpty(t, token, "Token should be present after registration")

	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})

	client := &http.Client{}
	checkResp, err := client.Do(req)
	require.NoError(t, err)
	defer checkResp.Body.Close()

	assert.Equal(t, http.StatusOK, checkResp.StatusCode)
}

func Test_CheckToken_MissingCookie(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/auth/check_token")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test_CheckToken_EmptyToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test_CheckToken_InvalidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "invalidtoken",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
