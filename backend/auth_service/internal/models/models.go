package models

type UserRegisterInfo struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	User_id uint   `json:"user_id"`
	Token   string `json:"token"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type CheckTokenRequest struct {
	BackendSessionToken string `json:"backend_session_token"`
}

type AdminLoginStruct struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
