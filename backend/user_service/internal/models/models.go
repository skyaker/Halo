package models

import "github.com/google/uuid"

type UserRegisterInfo struct {
	User_id  uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserDeleteInfo struct {
	User_id uuid.UUID `json:"id"`
}
