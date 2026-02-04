package service

import (
	entities "URL_checker/internal/repo/dto"
	"context"
	"time"
)

//Здесь будет основной сервис, делающий запросы и сохраняющий результат

type CheckRepository interface {
	Insert(ctx context.Context, r entities.Checks) error
	LatestByTarget(ctx context.Context, targetID int64) (entities.Checks, error)
	ListByTarget(ctx context.Context, targetID int64, limit int, from, to time.Time) ([]entities.Checks, error)
}
