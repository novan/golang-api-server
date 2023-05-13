package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Commander interface {
    MGet(ctx context.Context, client *redis.Client, keys ...string) ([]interface{}, error)
    Scan(data []interface{}, client *redis.Client, dst interface{}) error
    MSet(ctx context.Context, client *redis.Client, keyValues map[string]interface{}, duration time.Duration) error
    Del(ctx context.Context, client *redis.Client, keys ...string) error
}

type Cmd struct {
}

func NewCmd() Cmd {
	return Cmd{}
}

func (r *Cmd) MGet(ctx context.Context, client *redis.Client, keys ...string) ([]interface{}, error) {
	return client.MGet(ctx, keys...).Result()
}

func (r *Cmd) Scan(data []interface{}, dst interface{}) error {
	b, _ := json.Marshal(data)
	return json.Unmarshal(b, dst)
}

func (r *Cmd) MSet(ctx context.Context, client *redis.Client, keyValues map[string]interface{}, duration time.Duration) error {
	var pairs []interface{}
	var ttlCmd []interface{}
	ttlCmd = append(ttlCmd, "MULTI")
	for k, v := range keyValues {
		pairs = append(pairs, k, v)
		ttlCmd = append(ttlCmd, "EXPIRE", k, duration.Seconds())
	}
	ttlCmd = append(ttlCmd, "EXEC")

	if err := client.MSet(ctx, pairs).Err(); err != nil {
		return err
	}
	return client.Do(ctx, ttlCmd).Err()
}

func (r *Cmd) Del(ctx context.Context, client *redis.Client, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}
