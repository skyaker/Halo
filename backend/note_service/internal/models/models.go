package models

import (
	"github.com/google/uuid"
)

type NoteInfo struct {
	Id          uuid.UUID `json:"note_id"`
	Category_id uuid.UUID `json:"category_id"`
	Content     string    `json:"content"`
	Created_at  int64     `json:"created_at"`
	Updated_at  int64     `json:"updated_at"`
	Ended_at    int64     `json:"ended_at"`
	Completed   bool      `json:"completed"`
}

type NoteDeleteInfo struct {
	Note_id uuid.UUID `json:"note_id"`
}

type UserInfo struct {
	User_id uuid.UUID `json:"user_id"`
}

// type NoteGetInfo struct {
// 	User_id uuid.UUID ``
// }
