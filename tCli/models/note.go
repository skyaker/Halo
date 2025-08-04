package models

import (
	"github.com/google/uuid"
)

type NoteStruct struct {
	Id          uuid.UUID `json:"note_id"`
	Category_id uuid.UUID `json:"category_id"`
	Content     string    `json:"content"`
	Created_at  int       `json:"created_at"`
	Updated_at  int       `json:"updated_at"`
	Ended_at    int       `json:"ended_at"`
	Completed   bool      `json:"completed"`
	Synced      bool      `json:"synced"`
}
