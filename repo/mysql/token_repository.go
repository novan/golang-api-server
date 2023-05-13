package mysql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/novan/golang-api-server/repo/mysql/query"
	"github.com/novan/golang-api-server/repo/mysql/schema"
	"github.com/novan/golang-api-server/util"
	"github.com/pkg/errors"
)

const (
	createTokenError = "Gagal membuat token baru"
	updateTokenError = "Gagal menyimpan token baru"
)

type TokenRepositoryInterface interface {
	CreateToken(ctx context.Context, token string, jwtID string, userID int) error
	GetToken(ctx context.Context, refreshToken string) (*schema.UserToken, error)
	Update(ctx context.Context, token schema.UserToken) error
	GetFirstActive(ctx context.Context, userID int) (*schema.UserToken, error)
}

type TokenRepository struct {
	db    *sqlx.DB
	model MysqlInterface
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{
		db:    db,
		model: NewModel(),
	}
}

func (r *TokenRepository) CreateToken(ctx context.Context, token string, jwtID string, userID int) error {
	var expiredAt = time.Now().AddDate(0, 6, 0)
	var refreshToken = schema.UserToken{
		Token:       token,
		JwtID:       jwtID,
		CreatedAt:   time.Now(),
		ExpiredAt:   &expiredAt,
		IsUsed:      false,
		Invalidated: false,
		UserID:      userID,
	}

	_, err := r.model.PreparedCreateUpdate(ctx, r.db, refreshToken, query.USER_TOKENS_INSERT)
	if err != nil {
		return errors.Wrap(err, createTokenError)
	}
	return nil

}

func (r *TokenRepository) GetToken(ctx context.Context, refreshToken string) (*schema.UserToken, error) {
	var token schema.UserToken
	filter := map[string]interface{}{
		"token$eq": refreshToken,
	}
	err := r.model.Get(ctx, r.db, &token, util.NewQueryFilter(filter), query.USER_TOKENS_GET)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) Update(ctx context.Context, token schema.UserToken) error {
	_, err := r.model.PreparedCreateUpdate(ctx, r.db, token, query.USER_TOKENS_INSERT)
	if err != nil {
		return errors.Wrap(err, updateTokenError)
	}

	return nil
}

func (r *TokenRepository) GetFirstActive(ctx context.Context, userID int) (*schema.UserToken, error) {
	var token schema.UserToken
	// now := time.Now().Format(util.TIMEFORMAT_DATETIME)
	filter := map[string]interface{}{
		"user_id$eq!":   userID,
		"expired_at$gt": "NOW()",
	}
	err := r.model.Get(ctx, r.db, &token, util.NewQuery(0, 0, "expired_at DESC", filter, nil), query.USER_TOKENS_GET)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
