package auth_handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func extractUserIdFromToken(tokenStr string) (uuid.UUID, error) {
	secretKey, exists := os.LookupEnv("TOKEN_SECRET_KEY")
	if !exists {
		log.Error().Msg("Secret key not found")
		return uuid.UUID{}, fmt.Errorf("internal server error: secret key not found")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msg("Invalid token signing method")
			return nil, fmt.Errorf("bad request: invalid token signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Token parse error")
		return uuid.UUID{}, fmt.Errorf("bad request: token parse error")
	}

	claims, status := token.Claims.(jwt.MapClaims)
	if !status || !token.Valid {
		log.Error().Msg("Invalid token")
		return uuid.UUID{}, fmt.Errorf("unauthorized: invalid token")
	}

	userIdRaw, ok := claims["user_id"].(string)
	if !ok {
		log.Error().Msg("User id not found in token")
		return uuid.UUID{}, fmt.Errorf("unauthorized: user_id not found in token")
	}

	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		log.Error().Err(err).Msg("String to uuid parse error")
		return uuid.UUID{}, fmt.Errorf("internal server error: uuid parse failure")
	}

	log.Info().Str("process", "user id extraction from token").Msg("extraction succeed")

	return userId, nil
}

func CheckToken(db *sql.DB, redisDb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Error().Err(err).Msg("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user_session_token := session_cookie.Value

		if user_session_token == "" {
			log.Error().Msg("Session token is empty")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		user_id, err := extractUserIdFromToken(user_session_token)
		if err != nil {
			if strings.Contains(err.Error(), "unauthorized") {
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
			log.Error().Msg("User key doesn't exist")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokens, err := redisDb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min: minUnixTime,
			Max: maxUnixTime,
		}).Result()
		if err != nil {
			log.Error().Err(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// fmt.Print(tokens)

		tokenFound := slices.Contains(tokens, user_session_token)

		if !tokenFound {
			log.Error().Msg("Entered token is invalid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Info().Msg("Check token successful")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": user_id,
		})
	}
}
