package checks

import (
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
	"context"
	"time"
)

type ICheckHander interface {
	Insert(ctx context.Context, r dto.Checks) error
	LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error)
	ListByTarget(ctx context.Context, targetID uint64, limit int, from, to time.Time) ([]dto.Checks, error)
}

type CheckHandler struct {
	service service.ICheckService
}

func (h *CheckHandler) Insert(ctx context.Context, r dto.Checks) error {
	return h.service.Insert(ctx, r)
}

func (h *CheckHandler) LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error) {
	return h.service.LatestByTarget(ctx, targetID)
}

func (h *CheckHandler) ListByTarget(ctx context.Context, targetID uint64, limit int, from, to time.Time) ([]dto.Checks, error) {
	return h.service.ListByTarget(ctx, targetID, limit, from, to)
}
