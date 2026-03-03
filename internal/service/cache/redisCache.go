package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) (int64, error)
}

func NewCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		client: client,
	}
}

type RedisCache struct {
	client *redis.Client
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) (int64, error) {
	return c.client.Del(ctx, key).Result()
}
