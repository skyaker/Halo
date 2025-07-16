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

func TestDeleteUser_Success(t *testing.T) {
	userCode := uuid.NewString()[:8]
	user := map[string]string{
		"login":    "deluser_" + userCode,
		"password": "123456",
		"email":    "del" + userCode + "@example.com",
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

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/auth/delete_user", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: token})

	client := &http.Client{}
	delResp, err := client.Do(req)
	require.NoError(t, err)
	defer delResp.Body.Close()

	assert.Equal(t, http.StatusOK, delResp.StatusCode)
}

func TestDeleteUser_MissingToken(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/auth/delete_user", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteUser_TokenInvalid(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/auth/delete_user", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "fake.invalid.token",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteUser_AlreadyDeleted(t *testing.T) {
	userCode := uuid.NewString()[:8]
	user := map[string]string{
		"login":    "ghost_" + userCode,
		"password": "123456",
		"email":    "ghost" + userCode + "@example.com",
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

	req1, _ := http.NewRequest("DELETE", "http://localhost:8080/api/auth/delete_user", nil)
	req1.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	client := &http.Client{}
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	req2, _ := http.NewRequest("DELETE", "http://localhost:8080/api/auth/delete_user", nil)
	req2.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode) // или 404
}
