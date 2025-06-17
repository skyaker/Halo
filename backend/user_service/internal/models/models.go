package models

import "github.com/google/uuid"

type UserRegisterInfo struct {
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
