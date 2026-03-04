package serviceChecker

import (
	"URL_checker/internal/mapper"
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/service/cache"
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	DURATION = time.Minute * 5
)

type ICheckService interface {
	Insert(ctx context.Context, r dto.Checks) (dto.Checks, error)
	LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error)
	ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error)
}

type CheckService struct {
	repo  checks.ICheckRepository
	cache cache.IRedisCache
}

func NewCheckService(repo checks.ICheckRepository, cache cache.IRedisCache) ICheckService {
	return &CheckService{
		repo:  repo,
		cache: cache,
	}
}

func (s *CheckService) Insert(ctx context.Context, r dto.Checks) (dto.Checks, error) {
	key := strconv.Itoa(int(r.ID))
	value, err := mapper.FromCheck(r)
	if err != nil {
		fmt.Println("------------- NOT SETTING IN REDIS ---------------")
		return s.repo.Insert(ctx, r)
	}
	_ = s.cache.Set(ctx, key, value, DURATION)
	return s.repo.Insert(ctx, r)
}

func (s *CheckService) LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error) {
	return s.repo.LatestByTarget(ctx, targetID)
}

func (s *CheckService) ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error) {
	return s.repo.ListByTarget(ctx, targetID, limit)
}
