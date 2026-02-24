package configs

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	host     string
	Port     string
	user     string
	password string
	name     string
	sslMode  string
}

func (c *PostgresConfig) dsn() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.user, c.password),
		Host:   fmt.Sprintf("%s:%s", c.host, c.Port),
		Path:   c.name,
	}

	q := u.Query()
	q.Set("sslmode", c.sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}
func LoadPostgres() (*PostgresConfig, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	name := getEnv("DB_NAME", "postgres")
	sslMode := getEnv("DB_SSLMODE", "disable")

	if user == "" {
		return nil, errors.New("User must be set")
	}

	return &PostgresConfig{
		host:     host,
		Port:     port,
		user:     user,
		password: password,
		name:     name,
		sslMode:  sslMode,
	}, nil
}

func NewPostgres(cfg PostgresConfig) (*sql.DB, error) {
	var dsn string = cfg.dsn()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
