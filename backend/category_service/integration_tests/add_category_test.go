package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCategory_Success(t *testing.T) {
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

	category := map[string]interface{}{
		"name": "Test category",
	}
	categoryBody, _ := json.Marshal(category)

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
