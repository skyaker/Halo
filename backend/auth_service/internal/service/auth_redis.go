package auth_service

import (
	models "auth_service/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (s *authService) StartTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.deleteExpiredTokens(ctx)
			}
		}
	}()
}

func (s *authService) deleteExpiredTokens(ctx context.Context) {
	now := time.Now().Unix()

	users, err := s.redisDb.Keys(ctx, s.sessionPrefix+"*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching keys")
		return
	}

	for _, userKey := range users {
		removed, err := s.redisDb.ZRemRangeByScore(ctx, userKey, "0", fmt.Sprintf("%d", now)).
			Result()
		if err != nil {
			log.Error().Err(err).Msg("Error moving expired tokens")
			continue
		}

		if removed > 0 {
			log.Info().Msgf("Removed %d expired tokens from %s\n", removed, userKey)
		}
	}
}

func (s *authService) createToken(ctx context.Context, userId uuid.UUID) (string, error) {
	t := time.Now().Unix() + 60*60*24*31*6

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     t,
	})

	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}

	session_key := fmt.Sprintf("%v:%v", s.sessionPrefix, userId)

	_, err = s.redisDb.ZAdd(ctx, session_key, redis.Z{
		Score:  float64(t),
		Member: signed,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("token insertion failed: %w", err)
	}

	return signed, nil
}

func (s *authService) ParseToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt signing method not supported")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	claims, status := token.Claims.(jwt.MapClaims)
	if !status || !token.Valid {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	userIdRaw, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	return userId, nil
}

func (s *authService) deleteTokensByUserId(ctx context.Context, userId uuid.UUID) error {
	_, err := s.redisDb.Del(ctx, fmt.Sprintf("%v:%v", s.sessionPrefix, userId)).Result()
	if err != nil {
		return fmt.Errorf("redis token removal failed: %w", err)
	}
	return nil
}

func (s *authService) getUserTokens(ctx context.Context, userId uuid.UUID) ([]string, error) {
	key := fmt.Sprintf("%v:%v", s.sessionPrefix, userId)

	tokens, err := s.redisDb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", time.Now().Unix()),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("redis token fetch failed: %w", err)
	}

	if len(tokens) == 0 {
		return nil, models.ErrNotFound
	}

	return tokens, nil
}
