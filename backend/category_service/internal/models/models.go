package models

import "github.com/google/uuid"

type CategoryInfo struct {
	Id         uuid.UUID `json:"category_id"`
	User_id    uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	Created_at int64     `json:"created_at"`
	Updated_at int64     `json:"updated_at"`
}

type UserInfo struct {
	User_id uuid.UUID `json:"user_id"`
}
