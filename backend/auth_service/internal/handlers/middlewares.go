package auth_handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"net/http"
// 	"os"
// 	"slices"
// 	"strings"
// 	"time"
//
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/rs/zerolog/log"
// )
//
// func (h *AuthHandler) CheckToken() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session_cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			log.Error().Err(err).Msg("Unauthorized")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		user_session_token := session_cookie.Value
//
// 		if user_session_token == "" {
// 			log.Error().Msg("Session token is empty")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		user_id, errInfo := extractUserIdFromToken(user_session_token)
// 		if errInfo != nil {
// 			if strings.Contains(errInfo.Error(), "unauthorized") {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			} else if strings.Contains(errInfo.Error(), "bad request") {
// 				http.Error(w, "Bad request", http.StatusBadRequest)
// 			} else {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			}
// 			return
// 		}
//
// 		key := fmt.Sprintf("user:%v", user_id)
// 		minUnixTime := fmt.Sprintf("%d", time.Now().Unix())
// 		maxUnixTime := fmt.Sprintf("%d", math.MaxInt32-1)
//
// 		userExistence := h.redisDb.Get(context.Background(), key)
//
// 		if userExistence == nil {
// 			log.Error().Msg("User key doesn't exist")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		tokens, err := h.redisDb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
// 			Min: minUnixTime,
// 			Max: maxUnixTime,
// 		}).Result()
// 		if err != nil {
// 			log.Error().Err(err)
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
//
// 		tokenFound := slices.Contains(tokens, user_session_token)
//
// 		if !tokenFound {
// 			log.Error().Msg("Entered token is invalid")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		log.Info().Msg("Check token successful")
//
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"user_id": user_id,
// 		})
// 	}
// }
