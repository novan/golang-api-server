package schema

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/novan/golang-api-server/util"
)

type User struct {
	ID        int           `db:"id"`
	UserType  util.UserType `db:"user_type"`
	Email     string        `db:"email"`
	Password  string        `db:"password"`
	Mobile    *string       `db:"mobile"`
	LastLogin *time.Time    `db:"last_login"`
	IsActive  types.BitBool `db:"is_active"`
	CreatedBy int           `db:"created_by"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedBy int           `db:"updated_by"`
	UpdatedAt time.Time     `db:"updated_at"`
}

type UserToken struct {
	UserID      int           `db:"user_id"`
	Token       string        `db:"token"`
	JwtID       string        `db:"jwt_id"`
	CreatedAt   time.Time     `db:"created_at"`
	ExpiredAt   *time.Time    `db:"expired_at"`
	IsUsed      types.BitBool `db:"is_used"`
	Invalidated types.BitBool `db:"invalidated"`
}
