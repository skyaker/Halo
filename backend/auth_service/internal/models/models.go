package models

import "github.com/google/uuid"

type UserRegisterInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserCreatedEvent struct {
	User_id  uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserRegisterResponse struct {
	User_id uuid.UUID `json:"user_id"`
	Token   string    `json:"token"`
}

type UserDeletedEvent struct {
	Login string `json:"login"`
}

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type CheckTokenRequest struct {
	Token string `json:"token"`
}
