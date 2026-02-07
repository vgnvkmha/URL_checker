package target

import (
	entities "URL_checker/internal/repo/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type ITargetRepo interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	GetByID(ctx context.Context, id int) (entities.Targets, error)
	Update(ctx context.Context, id int, req entities.PatchReq) error
	Delete(ctx context.Context, id int) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
}

type TargetRepo struct {
	db *sql.DB
}

func New(db *sql.DB) (ITargetRepo, error) {
	return &TargetRepo{db: db}, nil
}

func (r *TargetRepo) Create(ctx context.Context, t entities.Targets) (entities.Targets, error) {

	query := `
		INSERT INTO targets (url, interval_sec, timeout_ms, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		t.URL,
		t.IntervalSec,
		t.TimeoutMS,
		t.Active,
	).Scan(
		&t.ID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return entities.Targets{}, err
	}

	return t, nil
}

func (r *TargetRepo) GetByID(ctx context.Context, id int) (entities.Targets, error) {

	var t entities.Targets

	query := `
		SELECT id, url, interval_sec, timeout_ms, active, created_at, updated_at
		FROM targets
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.URL,
		&t.IntervalSec,
		&t.TimeoutMS,
		&t.Active,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return entities.Targets{}, err
	}

	return t, nil
}

func (r *TargetRepo) Update(ctx context.Context, id int, req entities.PatchReq) error {

	setParts := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if req.Interval != nil {
		setParts = append(setParts, fmt.Sprintf("interval_sec = $%d", argID))
		args = append(args, *req.Interval)
		argID++
	}

	if req.Timeout != nil {
		setParts = append(setParts, fmt.Sprintf("timeout_ms = $%d", argID))
		args = append(args, *req.Timeout)
		argID++
	}

	if req.Active != nil {
		setParts = append(setParts, fmt.Sprintf("active = $%d", argID))
		args = append(args, *req.Active)
		argID++
	}

	if len(setParts) == 0 {
		return errors.New("Add params to edit")
	}

	query := fmt.Sprintf(`
		UPDATE targets
		SET %s,
		    updated_at = now()
		WHERE id = $%d
	`, strings.Join(setParts, ", "), argID)

	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TargetRepo) List(
	ctx context.Context,
) ([]entities.Targets, error) {
	query := `
		SELECT
			id,
			url,
			interval_sec,
			timeout_ms,
			active,
			created_at,
			updated_at
		FROM targets
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []entities.Targets

	for rows.Next() {
		var t entities.Targets

		if err := rows.Scan(
			&t.ID,
			&t.URL,
			&t.IntervalSec,
			&t.TimeoutMS,
			&t.Active,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		targets = append(targets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

func (r *TargetRepo) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM targets
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TargetRepo) ListActive(ctx context.Context) ([]entities.Targets, error) {
	query := `
		SELECT
			id,
			url,
			interval_sec,
			timeout_ms,
			active,
			created_at,
			updated_at
		FROM targets
		WHERE active=true
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []entities.Targets

	for rows.Next() {
		var t entities.Targets

		if err := rows.Scan(
			&t.ID,
			&t.URL,
			&t.IntervalSec,
			&t.TimeoutMS,
			&t.Active,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		targets = append(targets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
