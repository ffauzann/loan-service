package model

type LoginRequest struct {
	UserId     string `json:"user_id" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
	RememberMe bool   `json:"remember_me"`
}

type LoginResponse struct {
	Token    Token
	DeviceId string
}
