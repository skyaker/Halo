package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (s *authService) createToken(ctx context.Context, userId uuid.UUID) (string, error) {
	t := time.Now().Unix() + 60*60*24*31*6

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     t,
	})

	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("service: token creation failed: %w", err)
	}

	session_key := fmt.Sprintf("user:%v", userId)

	_, err = s.redisDb.ZAdd(ctx, session_key, redis.Z{
		Score:  float64(t),
		Member: signed,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("service: token insertion failed: %w", err)
	}

	return signed, nil
}
