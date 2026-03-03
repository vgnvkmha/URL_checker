package configs

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (c *RedisConfig) DSN() string {
	u := &url.URL{
		Scheme: "redis",
		Host:   fmt.Sprintf("%s:%s", c.Host, c.Port),
	}

	if c.Password != "" {
		u.User = url.UserPassword("", c.Password)
	}

	if c.DB != 0 {
		u.Path = "/" + string(rune(c.DB))
	}

	return u.String()
}

func Load() (*RedisConfig, error) {
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, err
	}

	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверка подключения
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
