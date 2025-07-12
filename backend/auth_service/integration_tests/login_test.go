package auth_integration_tests

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

func TestLogin_Success(t *testing.T) {
	user := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assertHasSessionToken(t, resp)
}

func TestLogin_InvalidLogin(t *testing.T) {
	user := map[string]string{
		"login":    "nonexistent_" + uuid.NewString()[:8],
		"password": "whatever",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLogin_InvalidPassword(t *testing.T) {
	user := map[string]string{
		"login":    "alice",
		"password": "wrongpass",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLogin_EmptyFields(t *testing.T) {
	user := map[string]string{
		"login":    "",
		"password": "",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLogin_InvalidJSON(t *testing.T) {
	invalidJSON := `{"login": "alice", "password": "alice123"`

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		strings.NewReader(invalidJSON),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLogin_SetsSessionTokenCookie(t *testing.T) {
	user := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assertHasSessionToken(t, resp)
}
