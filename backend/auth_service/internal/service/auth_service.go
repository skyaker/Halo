package auth_service

import (
	models "auth_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type AuthService interface {
	RegisterUser(ctx context.Context, userData models.UserRegisterInfo) (string, error)
	// DeleteUser() http.HandlerFunc
	// CheckToken() http.HandlerFunc
	// Login() http.HandlerFunc
}

type authService struct {
	db        *sql.DB
	redisDb   *redis.Client
	writer    *kafka.Writer
	secretKey string
}

func NewAuthService(
	db *sql.DB,
	redisDb *redis.Client,
	writer *kafka.Writer,
	key string,
) AuthService {
	return &authService{
		db:        db,
		redisDb:   redisDb,
		writer:    writer,
		secretKey: key,
	}
}

func (s *authService) RegisterUser(
	ctx context.Context,
	userData models.UserRegisterInfo,
) (string, error) {
	userId, err := addUserCredentials(ctx, s.db, &userData)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			return "", err
		}
		return "", fmt.Errorf("service: register user failed: %w", err)
	}

	createdEvent := models.UserCreatedEvent{
		User_id:  userId,
		Username: userData.Username,
		Email:    userData.Email,
	}

	err = s.sentUserCreatedEvent(ctx, userId, createdEvent)
	if err != nil {
		return "", fmt.Errorf("service: send user created event failed: %w", err)
	}

	token, err := s.createToken(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("service: token creation failed: %w", err)
	}
	return token, nil
}

func (s *authService) DeleteUser() {
}

// func (s *authService) CheckToken() http.HandlerFunc {
// }
//
// func (s *authService) Login() http.HandlerFunc {
// }
