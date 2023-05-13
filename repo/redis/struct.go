package redis

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/novan/golang-api-server/util"
)

type User struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	UserType  util.UserType `json:"user_type"`
	Email     string        `json:"email"`
	Mobile    *string       `json:"mobile"`
	LastLogin *time.Time    `json:"last_login"`
	IsActive  types.BitBool `json:"is_active"`
}
