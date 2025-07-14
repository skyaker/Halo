package auth_integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func GetSessionTokenCookie(t *testing.T, login, password string) string {
	user := map[string]string{
		"login":    login,
		"password": password,
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	for _, c := range resp.Cookies() {
		if c.Name == "session_token" {
			return c.Value
		}
	}

	t.Fatal("session_token not found in login response")
	return ""
}

func Test_CheckToken_Success(t *testing.T) {
	token := GetSessionTokenCookie(t, "alice", "alice123")

	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
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
}

func Test_CheckToken_MissingCookie(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/auth/check_token")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test_CheckToken_EmptyToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
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
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/auth/check_token", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "thisisnotavalidtoken",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
