package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepositoryInterface interface {
	StoreUser(ctx context.Context, token string, user *User) error 
	GetUser(ctx context.Context, token string) (user *User, err error)
	RemoveUser(ctx context.Context, token string) (err error)
}

type SessionRepository struct {
	dbRedis *redis.Client
}

func NewSessionRepository(r *redis.Client) *SessionRepository {
	return &SessionRepository{dbRedis: r}
}

func (sm *SessionRepository) StoreUser(ctx context.Context, token string, user *User) error {
	obj, err := json.Marshal(user)
    if err != nil {
       return err
    }
	return sm.dbRedis.Set(ctx, token, obj, 30 * 24 * time.Hour).Err()
}

func (sm *SessionRepository) GetUser(ctx context.Context, token string) (user *User, err error) {
	data, err := sm.dbRedis.Get(ctx, token).Result()
	if err == nil {
		err = json.Unmarshal([]byte(data), &user)
	}
	return
}

func (sm *SessionRepository) RemoveUser(ctx context.Context, token string) (err error) {
	_, err = sm.dbRedis.Del(ctx, token).Result()
	return
}