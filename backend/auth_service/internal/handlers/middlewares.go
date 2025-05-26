package auth_handlers

import (
	models "auth_service/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func checkBackendToken(tokenStr *string) error {
	backend_session_token, exists := os.LookupEnv("BACKEND_SESSION_TOKEN")
	if !exists {
		return fmt.Errorf("internal server error: backend session token was not found")
	}

	if *tokenStr != backend_session_token {
		return fmt.Errorf("declined: backend session token was not confirmed")
	}

	return fmt.Errorf("success: backend session token was confirmed")
}

func extractUserIdFromToken(tokenStr string) (uint, error) {
	secretKey, exists := os.LookupEnv("TOKEN_SECRET_KEY")
	if !exists {
		return 0, fmt.Errorf("internal server error: secret key not found")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad request: invalid token signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("bad request: token parse error")
	}

	claims, status := token.Claims.(jwt.MapClaims)
	if !status || !token.Valid {
		return 0, fmt.Errorf("unauthorized: invalid token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unauthorized: user_id not found in token")
	}

	return uint(userIDFloat), nil
}

func CheckToken(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkTokenRequest models.CheckTokenRequest

		if err := json.NewDecoder(r.Body).Decode(&checkTokenRequest); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		session_cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user_session_token := session_cookie.Value
		backend_session_token := checkTokenRequest.BackendSessionToken
		fmt.Print(backend_session_token)

		if user_session_token == "" || backend_session_token == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		backendTokenErr := checkBackendToken(&backend_session_token)
		if !strings.Contains(backendTokenErr.Error(), "success") {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		user_id, err := extractUserIdFromToken(user_session_token)
		if err != nil {
			if strings.Contains(err.Error(), "unauthorized") {
				// log (?)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else if strings.Contains(err.Error(), "bad request") {
				http.Error(w, "Bad request", http.StatusBadRequest)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		key := fmt.Sprintf("user:%v", user_id)
		minUnixTime := fmt.Sprintf("%d", time.Now().Unix())
		maxUnixTime := fmt.Sprintf("%d", math.MaxInt32-1)

		userExistence := redisDb.Get(context.Background(), key)

		if userExistence == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokens, err := redisDb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min: minUnixTime,
			Max: maxUnixTime,
		}).Result()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenFound := false
		for _, storedToken := range tokens {
			if storedToken == user_session_token {
				tokenFound = true
				break
			}
		}

		if !tokenFound {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": user_id,
		})
	}
}

func VerifyAdminToken(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminCookie, err := r.Cookie("admin_token")
		if err != nil {
			http.Error(w, "Unauthorized: admin token missing", http.StatusUnauthorized)
			return
		}

		adminToken := adminCookie.Value
		if adminToken == "" {
			http.Error(w, "Unauthorized: empty admin token", http.StatusUnauthorized)
			return
		}

		adminID, err := extractAdminIdFromToken(adminToken)
		if err != nil {
			if strings.Contains(err.Error(), "unauthorized") {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		key := fmt.Sprintf("admin:%v", adminID)
		minUnixTime := fmt.Sprintf("%d", time.Now().Unix())
		maxUnixTime := fmt.Sprintf("%d", math.MaxInt32-1)

		tokens, err := redisDb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min: minUnixTime,
			Max: maxUnixTime,
		}).Result()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tokenFound := false
		for _, storedToken := range tokens {
			if storedToken == adminToken {
				tokenFound = true
				break
			}
		}

		if !tokenFound {
			http.Error(w, "Unauthorized: token not found in Redis", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"admin_id": adminID,
		})
	}
}

func extractAdminIdFromToken(tokenStr string) (uint, error) {
	secretKey, exists := os.LookupEnv("ADMIN_TOKEN_SECRET_KEY")
	if !exists {
		return 0, fmt.Errorf("internal server error: admin secret key not found")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unauthorized: invalid token signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("unauthorized: token parse error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("unauthorized: invalid token claims")
	}

	adminIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unauthorized: admin_id not found in token")
	}

	return uint(adminIDFloat), nil
}
