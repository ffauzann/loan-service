package model

type RegisterRequestUserDetail struct {
	Name string `json:"name" validate:"required"`

	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	RoleId      uint8  `json:"role_id"`

	PlainPassword string `json:"password" name:"password" validate:"required,password"`
	UserPassword  string
}

type RegisterRequestGroupDetail struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	AllowSupervision bool   `json:"allow_supervision"`
}

type RegisterRequest struct {
	User  RegisterRequestUserDetail   `json:"user"`
	Group *RegisterRequestGroupDetail `json:"group"`
}

type RegisterResponse struct {
	StatusCode int
	Reasons    []string
}
