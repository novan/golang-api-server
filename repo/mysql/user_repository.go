package mysql

import (
	"context"

	"github.com/novan/golang-api-server/repo/mysql/query"
	"github.com/novan/golang-api-server/repo/mysql/schema"
	"github.com/novan/golang-api-server/util"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	FindByID(ctx context.Context, userID int) (*schema.User, error)
	FindByEmail(ctx context.Context, email string) (*schema.User, error)
	FindByToken(ctx context.Context, token string) (*schema.User, error)
	Create(ctx context.Context, user schema.User) (int, error)
	Update(ctx context.Context, user schema.User) error
	Login(ctx context.Context, email string, password string) (*schema.User, error)
	UpdatePassword(ctx context.Context, ID int, password string) error
}

type UserRepository struct {
	db    *sqlx.DB
	model MysqlInterface
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:    db,
		model: NewModel(),
	}
}

func (ar *UserRepository) FindByID(ctx context.Context, userID int) (*schema.User, error) {
	var user schema.User
	filter := map[string]interface{}{
		"id$eq": userID,
	}

	err := ar.model.Get(ctx, ar.db, &user, util.NewQueryFilter(filter), query.USER_LIST)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (ar *UserRepository) FindByEmail(ctx context.Context, email string) (*schema.User, error) {
	var user schema.User
	filter := map[string]interface{}{
		"Email$eq": email,
	}

	err := ar.model.Get(ctx, ar.db, &user, util.NewQueryFilter(filter), query.USER_LIST)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (ar *UserRepository) FindByToken(ctx context.Context, token string) (*schema.User, error) {
	var user schema.User
	filter := map[string]interface{}{
		"ut.token$eq": token,
	}

	err := ar.model.Get(ctx, ar.db, &user, util.NewQueryFilter(filter), query.USER_TOKEN)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (ar *UserRepository) Create(ctx context.Context, user schema.User) (int, error) {
	result, err := ar.model.CreateUpdate(ctx, ar.db, user, query.USER_CREATE)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (ar *UserRepository) Update(ctx context.Context, user schema.User) error {
	_, err := ar.model.CreateUpdate(ctx, ar.db, user, query.USER_UPDATE)
	if err != nil {
		return err
	}
	return nil
}

func (ar *UserRepository) Login(ctx context.Context, email string, password string) (*schema.User, error) {
	var user schema.User
	filter := map[string]interface{}{
		"email$eq":    email,
		"password$eq": password,
	}

	err := ar.model.Get(ctx, ar.db, &user, util.NewQueryFilter(filter), query.USER_LIST)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ar *UserRepository) UpdatePassword(ctx context.Context, ID int, password string) error {
	var user = schema.User{
		ID:       ID,
		Password: password,
	}

	_, err := ar.model.CreateUpdate(ctx, ar.db, user, query.USER_UPDATE_PASSWORD)
	return err
}
