package models

import "github.com/google/uuid"

type NoteInfo struct {
	User_id uuid.UUID `json:"user_id"`
	Type_id uuid.UUID `json:"type_id"`
	Content string    `json:"content"`
}
