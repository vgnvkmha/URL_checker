package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) (int64, error)
	GetAll(ctx context.Context) (map[string]interface{}, error)
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

func (c *RedisCache) GetAll(ctx context.Context) (map[string]interface{}, error) {
	var cursor uint64
	result := make(map[string]interface{})

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return nil, err
		}

		if len(keys) > 0 {
			values, err := c.client.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, err
			}

			for i, key := range keys {
				if values[i] == nil {
					continue
				}

				strVal := values[i].(string)

				var parsed interface{}
				if err := json.Unmarshal([]byte(strVal), &parsed); err != nil {
					// если вдруг не JSON — оставляем строку
					result[key] = strVal
					continue
				}

				result[key] = parsed
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}
