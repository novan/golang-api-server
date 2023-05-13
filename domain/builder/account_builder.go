package builder

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/novan/golang-api-server/domain/entity"
	"github.com/novan/golang-api-server/repo/mysql/schema"
	repo "github.com/novan/golang-api-server/repo/redis"
	httpModel "github.com/novan/golang-api-server/transport/http/model"
	"github.com/novan/golang-api-server/util"
)

type AccountBuilder struct{}

func NewAccountBuilder() *AccountBuilder {
	return &AccountBuilder{}
}

func (b *AccountBuilder) ToDBModel(user *entity.User) *schema.User {
	return &schema.User{
		UserType:  user.UserType,
		Email:     user.Email,
		Mobile:    user.Mobile,
		LastLogin: user.LastLogin,
		IsActive:  types.BitBool(user.IsActive),
	}
}

func (b *AccountBuilder) ToDomainModel(user *schema.User) *entity.User {
	return &entity.User{
		ID:        user.ID,
		UserType:  user.UserType,
		Email:     user.Email,
		Mobile:    user.Mobile,
		LastLogin: user.LastLogin,
		IsActive:  bool(user.IsActive),
	}
}

func (b *AccountBuilder) FromRepoToSessionModel(user *schema.User) *repo.User {
	return &repo.User{
		ID:        user.ID,
		UserType:  util.UserType(user.UserType),
		Email:     user.Email,
		Mobile:    user.Mobile,
		LastLogin: user.LastLogin,
		IsActive:  user.IsActive,
	}
}

func (b *AccountBuilder) FromSignupRequestToDomain(user *httpModel.SignupRequest) entity.SignupRequest {
	return entity.SignupRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Mobile:   user.Mobile,
		UserType: user.UserType,
	}
}
