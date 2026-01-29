package models

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
