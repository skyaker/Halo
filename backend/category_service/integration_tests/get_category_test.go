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

func TestGetCategory_Success(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken := loginAndGetToken(t, creds)

	nameId := uuid.New()
	name := fmt.Sprintf("Category name %s", nameId.String())

	category := map[string]interface{}{
		"name": name,
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

	getReq, err := http.NewRequest(
		"GET",
		"http://localhost:8080/api/category?page=1",
		nil,
	)
	getReq.AddCookie(&http.Cookie{Name: "session_token", Value: sessionToken})
	getResp, err := client.Do(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()

	var categories []map[string]interface{}
	err = json.NewDecoder(getResp.Body).Decode(&categories)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
}

func TestGetCategory_NoToken(t *testing.T) {
	req, err := http.NewRequest(
		"GET",
		"http://localhost:8080/api/category?page=1",
		nil,
	)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetCategory_InvalidToken(t *testing.T) {
	req, err := http.NewRequest(
		"GET",
		"http://localhost:8080/api/category?page=1",
		nil,
	)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "invalid"})
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetCategory_LastPage(t *testing.T) {
	creds := map[string]string{
		"login":    "alice",
		"password": "alice123",
	}

	sessionToken := loginAndGetToken(t, creds)

	client := &http.Client{}

	getReq, err := http.NewRequest(
		"GET",
		"http://localhost:8080/api/category?page=999",
		nil,
	)
	getReq.AddCookie(&http.Cookie{Name: "session_token", Value: sessionToken})
	getResp, err := client.Do(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()

	var categories []map[string]interface{}
	err = json.NewDecoder(getResp.Body).Decode(&categories)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
}

