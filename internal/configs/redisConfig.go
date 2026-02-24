package configs

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
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

	// Если есть пароль
	if c.Password != "" {
		u.User = url.UserPassword("", c.Password)
	}

	// Номер базы (db)
	if c.DB != 0 {
		u.Path = "/" + string(c.DB)
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
