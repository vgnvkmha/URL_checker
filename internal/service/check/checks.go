package serviceChecker

import (
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"context"

	"github.com/redis/go-redis/v9"
)

type ICheckService interface {
	Insert(ctx context.Context, r dto.Checks) (dto.Checks, error)
	LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error)
	ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error)
}

type CheckService struct {
	repo  checks.ICheckRepository
	redis *redis.Client
}

func NewCheckService(repo checks.ICheckRepository, redisClient *redis.Client) ICheckService {
	return &CheckService{
		repo:  repo,
		redis: redisClient,
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
