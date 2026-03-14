package serviceChecker

import (
	"URL_checker/internal/mapper"
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/service/cache"
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"
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
	repo   checks.ICheckRepository
	cache  cache.IRedisCache
	logger *zap.SugaredLogger
}

func NewCheckService(repo checks.ICheckRepository, cache cache.IRedisCache, logger *zap.SugaredLogger) ICheckService {
	return &CheckService{
		repo:  repo,
		cache: cache,
		logger: logger.With(
			"module", "check",
			"layer", "service"),
	}
}

func (s *CheckService) Insert(ctx context.Context, r dto.Checks) (dto.Checks, error) {
	key := strconv.Itoa(int(r.ID))
	const operation = "Insert"
	value, err := mapper.FromCheck(r)
	if err != nil {
		s.logger.Errorw("unsuccesful setting in redis error",
			"error", err)
		return s.repo.Insert(ctx, r)
	}
	_ = s.cache.Set(ctx, key, value, DURATION)
	s.logger.Infow("Succesful",
		"operation", operation,
	)
	return s.repo.Insert(ctx, r)
}

func (s *CheckService) LatestByTarget(ctx context.Context, targetID uint64) (dto.Checks, error) {
	target, err := s.repo.LatestByTarget(ctx, targetID)
	if err != nil {
		s.logger.Errorw("postgres error",
			"error", err)
		return dto.Checks{}, err
	}
	return target, nil
}

func (s *CheckService) ListByTarget(ctx context.Context, targetID uint64, limit int) ([]dto.Checks, error) {
	target, err := s.repo.ListByTarget(ctx, targetID, limit)
	if err != nil {
		s.logger.Errorw("postgres error",
			"error", err)
		return []dto.Checks{}, err
	}
	return target, nil
}
