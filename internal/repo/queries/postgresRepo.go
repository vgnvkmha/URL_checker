package queries

import (
	entities "URL_checker/internal/repo/dto"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type TargetRepo interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	GetByID(ctx context.Context, id int) (entities.Targets, error)
	Update(ctx context.Context, t entities.Targets) error
	Delete(ctx context.Context, id int) error
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(dsn string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err // не смогли открыть соединение
	}

	if err := db.Ping(); err != nil {
		return nil, err // база не отвечает
	}

	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) Create(
	ctx context.Context,
	t entities.Targets,
) (entities.Targets, error) {

	query := `
		INSERT INTO targets1 (url, interval_sec, timeout_ms, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		t.URL,
		t.IntervalSec,
		t.TimeoutMS,
		t.Active,
	).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		return entities.Targets{}, err
	}

	return t, nil
}
