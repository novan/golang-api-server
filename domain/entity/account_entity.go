package entity

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/novan/golang-api-server/util"
)

type User struct {
	ID        int           `json:"id"`
	UserType  util.UserType `json:"user_type"`
	Email     string        `json:"email"`
	Mobile    *string       `json:"mobile"`
	LastLogin *time.Time    `json:"last_login"`
	IsActive  bool          `json:"is_active"`
}

type UserToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type CustomClaims struct {
	UserID   int    `json:"id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

type SignupRequest struct {
	Name     string
	Email    string
	Password string
	Mobile   *string
	UserType util.UserType
}
