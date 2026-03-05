package auth_service

import (
	models "auth_service/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func (s *authService) sentUserCreatedEvent(
	ctx context.Context,
	userId uuid.UUID,
	userData models.UserCreatedEvent,
) error {
	event := map[string]any{
		"id":       userId,
		"username": "",
		"email":    "",
	}

	if userData.Username != "" {
		event["username"] = userData.Username
	}

	if userData.Email != "" {
		event["email"] = userData.Email
	}

	data, _ := json.Marshal(event)

	err := s.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user-created",
		Key:   []byte(userId.String()),
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("kafka user-created message error: %w", err)
	}

	return nil
}

func (s *authService) sentUserDeletedEvent(
	ctx context.Context,
	userId uuid.UUID,
) error {
	err := s.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user-deleted",
		Key:   []byte(userId.String()),
		Value: []byte(userId.String()),
	})
	if err != nil {
		return fmt.Errorf("kafka user-deleted message error: %w", err)
	}
	return nil
}
