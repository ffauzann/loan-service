package model

import "github.com/ffauzann/loan-service/internal/constant"

type User struct {
	CommonModel

	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`

	Password []byte `json:"password" db:"password"`

	Status          constant.UserStatus `json:"status" db:"status"`
	IsEmailVerified bool                `json:"is_email_verified" db:"is_email_verified"`
	RoleId          uint8               `json:"role_id" db:"role_id"`
}

type IsUserExistRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type IsUserExistResponse struct {
	IsExist bool     `json:"is_exist"`
	Reasons []string `json:"reasons"`
}

type CloseAccountRequest struct {
	UserId uint64 `json:"-"` // From claims.
}

type CloseAccountResponse struct {
	UserId uint64              `json:"user_id"`
	Status constant.UserStatus `json:"status"`
}
