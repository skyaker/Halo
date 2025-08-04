package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loginAndGetToken(t *testing.T, creds map[string]string) string {
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

	return sessionToken
}

func TestAddCategory_Success(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken := loginAndGetToken(t, creds)

	nameId := uuid.New()

	category := map[string]interface{}{
		"name": fmt.Sprintf("Test category %s", nameId.String()),
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

func TestAddCategory_EmptyName(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken := loginAndGetToken(t, creds)

	category := map[string]interface{}{
		"name": "",
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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAddCategory_NoToken(t *testing.T) {
	nameId := uuid.New()

	category := map[string]interface{}{
		"name": fmt.Sprintf("Unauthorized %s", nameId.String()),
	}
	categoryBody, _ := json.Marshal(category)

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAddCategory_DuplicateName(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken := loginAndGetToken(t, creds)

	nameId := uuid.New()
	name := fmt.Sprintf("Repeated category name %s", nameId.String())

	category := map[string]interface{}{
		"name": name,
	}
	categoryBody, _ := json.Marshal(category)

	client := &http.Client{}

	req1, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
	)
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	req2, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
	)
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
}

func TestAddCategory_DuplicateNameWithDifferentUser(t *testing.T) {
	creds1 := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken1 := loginAndGetToken(t, creds1)

	creds2 := map[string]string{
		"login":    "bob",
		"password": "bob123",
	}

	sessionToken2 := loginAndGetToken(t, creds2)

	nameId := uuid.New()
	name := fmt.Sprintf("Repeated category name %s", nameId.String())

	category := map[string]interface{}{
		"name": name,
	}
	categoryBody, _ := json.Marshal(category)

	client := &http.Client{}

	req1, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
	)
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")
	req1.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken1,
	})
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	req2, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/category",
		bytes.NewBuffer(categoryBody),
	)
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken2,
	})
	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}
