package auth_service

import (
	models "auth_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type AuthService interface {
	RegisterUser(ctx context.Context, userData models.UserRegisterInfo) (string, error)
	DeleteUser(ctx context.Context, userSessionToken string) error
	CheckToken(ctx context.Context, userSessionToken string) (uuid.UUID, error)
	Login(ctx context.Context, userData models.UserLogin) (string, error)
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
	if userData.Login == "" || userData.Password == "" {
		return "", models.ErrInvalidRequest
	}

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

func (s *authService) DeleteUser(ctx context.Context, userSessionToken string) error {
	if userSessionToken == "" {
		return models.ErrInvalidToken
	}

	userId, err := s.ParseToken(userSessionToken)
	if err != nil {
		return fmt.Errorf("service: token parse failed: %w", err)
	}

	err = deleteUserCredentials(ctx, s.db, userId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.ErrNotFound
		}
		return fmt.Errorf("service: delete user credentials failed: %w", err)
	}

	err = s.sentUserDeletedEvent(ctx, userId)
	if err != nil {
		return fmt.Errorf("service: send user deleted event failed: %w", err)
	}

	err = s.deleteTokensByUserId(ctx, userId)
	if err != nil {
		return fmt.Errorf("service: delete tokens by user id failed: %w", err)
	}

	return nil
}

func (s *authService) CheckToken(ctx context.Context, userSessionToken string) (uuid.UUID, error) {
	if userSessionToken == "" {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	userId, err := s.ParseToken(userSessionToken)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service: token parse failed: %w", err)
	}

	tokens, err := s.getUserTokens(ctx, userId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service: redis token fetch failed: %w", err)
	}

	tokenFound := slices.Contains(tokens, userSessionToken)
	if !tokenFound {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	return userId, nil
}

func (s *authService) Login(ctx context.Context, userData models.UserLogin) (string, error) {
	if userData.Login == "" || userData.Password == "" {
		return "", models.ErrInvalidRequest
	}

	userId, err := checkCredentials(ctx, s.db, userData.Login, userData.Password)
	if err != nil {
		return "", fmt.Errorf("service: check credentials failed: %w", err)
	}

	token, err := s.createToken(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("service: token creation failed: %w", err)
	}

	return token, nil
}
