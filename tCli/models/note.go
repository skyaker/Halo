package models

import (
	"time"

	"github.com/google/uuid"
)

type NoteStruct struct {
	Id         uuid.UUID `json:"note_id"`
	Type_id    uuid.UUID `json:"type_id"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Ended_at   time.Time `json:"ended_at"`
	Completed  bool      `json:"completed"`
	Synced     bool      `json:"synced"`
}
