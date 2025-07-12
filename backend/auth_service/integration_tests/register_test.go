package auth_integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertHasSessionToken(t *testing.T, resp *http.Response) {
	var sessionFound bool
	for _, c := range resp.Cookies() {
		if c.Name == "session_token" {
			sessionFound = true
			break
		}
	}
	assert.True(t, sessionFound, "session_token cookie must be set")
}

func TestRegister_Success(t *testing.T) {
	userCode := uuid.NewString()[:8]
	user := map[string]string{
		"login":    "testuser_" + userCode,
		"password": "1234",
		"email":    "test" + userCode + "@example.com",
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assertHasSessionToken(t, resp)
}

func TestRegister_DuplicateLogin(t *testing.T) {
	user := map[string]string{
		"login":    "alice",
		"password": "anypass",
		"username": "AliceClone",
		"email":    "alice-clone@example.com",
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegister_EmptyFields(t *testing.T) {
	user := map[string]string{
		"login":    "",
		"password": "",
		"username": "",
		"email":    "",
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegister_InvalidJSON(t *testing.T) {
	badJSON := `{"login": "broken", password: "1234"}`

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBufferString(badJSON),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegister_ExtraFields(t *testing.T) {
	userCode := uuid.NewString()[:8]

	user := map[string]any{
		"login":         "extra_" + userCode,
		"password":      "123456",
		"username":      "extra",
		"email":         "extra" + userCode + "@example.com",
		"unexpected":    "ex",
		"another_field": 123,
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assertHasSessionToken(t, resp)
}

func TestRegister_TooLongFields(t *testing.T) {
	longStr := strings.Repeat("x", 300)

	user := map[string]string{
		"login":    longStr,
		"password": "secure123",
		"username": "TooLongUser",
		"email":    longStr + "@example.com",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Contains(
		t,
		[]int{http.StatusBadRequest, http.StatusInternalServerError},
		resp.StatusCode,
	)
}

func TestRegister_NoOptionalFields(t *testing.T) {
	userCode := uuid.NewString()[:8]

	user := map[string]string{
		"login":    "onlylogin_" + userCode,
		"password": "onlypass",
	}

	body, _ := json.Marshal(user)

	resp, err := http.Post(
		"http://localhost:8080/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assertHasSessionToken(t, resp)
}
