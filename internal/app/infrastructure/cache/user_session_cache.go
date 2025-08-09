package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserSessionCache struct {
	client *redis.Client
}

func NewUserSessionCache(client *redis.Client) *UserSessionCache {
	return &UserSessionCache{client}
}

func (usc *UserSessionCache) SaveToken(ctx context.Context, token string, ttl time.Duration) error {
	return usc.client.SetArgs(ctx, token, "", redis.SetArgs{TTL: ttl, Mode: "NX"}).Err()
}

func (usc *UserSessionCache) CheckToken(ctx context.Context, token string) (bool, error) {
	err := usc.client.Get(ctx, token).Err()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (usc *UserSessionCache) RevokeToken(ctx context.Context, token string) error {
	return usc.client.Del(ctx, token).Err()
}
