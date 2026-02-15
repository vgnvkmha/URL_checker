package serviceChecker

import (
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"context"
)

type ICheckService interface {
	Insert(ctx context.Context, r dto.Checks) (dto.Checks, error)
	LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error)
	ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error)
}

type CheckService struct {
	repo checks.ICheckRepository
}

func NewCheckService(repo checks.ICheckRepository) ICheckService {
	return &CheckService{
		repo: repo,
	}
}

func (s *CheckService) Insert(ctx context.Context, r dto.Checks) (dto.Checks, error) {
	return s.repo.Insert(ctx, r)
}

func (s *CheckService) LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error) {
	return s.repo.LatestByTarget(ctx, targetID)
}

func (s *CheckService) ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error) {
	return s.repo.ListByTarget(ctx, targetID, limit)
}
