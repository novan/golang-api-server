package model

import (
	"github.com/novan/golang-api-server/util"
)

type Token struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	Fullname string `json:"fullname"`
	Mobile   string `json:"mobile"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Name     string        `json:"name" validate:"required"`
	Email    string        `json:"email" validate:"required,email"`
	Password string        `json:"password" validate:"required"`
	Mobile   *string       `json:"mobile"`
	UserType util.UserType `json:"user_type" validate:"omitempty,oneof=ADMIN USER"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ToggleActiveRequest struct {
	UserID   int  `json:"user_id" validate:"required,numeric"`
	IsActive bool `json:"is_active" validate:"required"`
}
