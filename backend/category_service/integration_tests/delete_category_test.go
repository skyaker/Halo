package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	models "category_service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCategory_Success(t *testing.T) {
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
		"http://localhost:8080/api/note?page=1",
		nil,
	)

	getReq.Header.Set("Content-Type", "application/json")
	getReq.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	getResp, err := client.Do(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()

	var notes []models.CategoryInfo
	err = json.NewDecoder(getResp.Body).Decode(&notes)
	require.NoError(t, err)

	deleteReq, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://localhost:8080/api/category?id=%s", nameId.String()),
		nil,
	)
	require.NoError(t, err)

	deleteReq.Header.Set("Content-Type", "application/json")
	deleteReq.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})

	deleteResp, err := client.Do(deleteReq)
	require.NoError(t, err)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
}
