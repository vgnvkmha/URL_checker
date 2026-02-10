package checks

import (
	"URL_checker/internal/repo/dto"
	"context"
	"database/sql"
	"log"
	"time"
)

type ICheckRepository interface {
	Insert(ctx context.Context, r dto.Checks) (dto.Checks, error)
	LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error)
	ListByTarget(ctx context.Context, targetID uint64, limit int, from, to time.Time) ([]dto.Checks, error)
}
type CheckRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (*CheckRepository, error) {
	return &CheckRepository{db: db}, nil
}

func (rCheck *CheckRepository) Insert(ctx context.Context, r dto.Checks) (dto.Checks, error) {
	row := rCheck.db.QueryRowContext(ctx, `
    SELECT
        current_database(),
        current_schema(),
        current_setting('search_path')
	`)
	var db, schema, path string
	_ = row.Scan(&db, &schema, &path)

	log.Printf("DB=%s SCHEMA=%s SEARCH_PATH=%s", db, schema, path)

	query := `
		INSERT INTO checks (target_id, ok, status_code, latency_ms, error)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, checked_at
	`
	err := rCheck.db.QueryRowContext(
		ctx,
		query,
		r.TargetId,
		r.OK,
		r.StatusCode,
		r.LatencyMs,
		r.Error,
	).Scan(
		&r.ID,
		&r.CheckedAt,
	)
	if err != nil {
		return dto.Checks{}, err
	}

	return r, nil
}

func (rCheck *CheckRepository) LatestByTarget(
	ctx context.Context,
	targetID uint64,
) (dto.Checks, error) {

	query := `
		SELECT *
		FROM checks
		WHERE target_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`

	var check dto.Checks

	err := rCheck.db.QueryRowContext(
		ctx,
		query,
		targetID,
	).Scan(
		&check.ID,
		&check.TargetId,
		&check.CheckedAt,
		&check.OK,
		&check.StatusCode,
		&check.LatencyMs,
		&check.Error,
	)

	if err != nil {
		return dto.Checks{}, err
	}

	return check, nil
}

func (rCheck *CheckRepository) ListByTarget(
	ctx context.Context,
	targetID uint64,
	limit int,
	from, to time.Time,
) ([]dto.Checks, error) {

	query := `
		SELECT
			id,
			target_id,
			checked_at,
			ok,
			status_code,
			latency_ms,
			error
		FROM checks
		WHERE target_id = $1
		  AND checked_at BETWEEN $2 AND $3
		ORDER BY checked_at DESC
		LIMIT $4
	`

	rows, err := rCheck.db.QueryContext(
		ctx,
		query,
		targetID,
		from,
		to,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []dto.Checks

	for rows.Next() {
		var c dto.Checks

		if err := rows.Scan(
			&c.ID,
			&c.TargetId,
			&c.CheckedAt,
			&c.OK,
			&c.StatusCode,
			&c.LatencyMs,
			&c.Error,
		); err != nil {
			return nil, err
		}

		checks = append(checks, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checks, nil
}
