package serviceTarget

import (
	"URL_checker/internal/mapper"
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"URL_checker/internal/service/cache"
	"URL_checker/internal/service/validation"
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

const (
	DURATION = time.Hour
)

type IURLService interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id int) (entities.Targets, error)
	List(ctx context.Context) (map[string]interface{}, error)
	Update(ctx context.Context, id int, params entities.PatchReq) error
	Delete(ctx context.Context, id int) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type URLService struct {
	repo   target.ITargetRepo
	cache  cache.IRedisCache
	group  *singleflight.Group
	logger *zap.SugaredLogger
}

func New(repo target.ITargetRepo, cache cache.IRedisCache,
	group *singleflight.Group, logger *zap.SugaredLogger) *URLService { //TODO: посмотреть, на что ругается
	return &URLService{
		repo:  repo,
		cache: cache,
		group: group,
		logger: logger.With(
			"module", "target",
			"layer", "service",
		),
	}
}

func (s *URLService) Create(
	ctx context.Context,
	t entities.Targets,
) (entities.Targets, error) {

	const operation = "create"
	err := validation.ValidURL(t.URL)
	if err != nil {
		s.logger.Errorw("Invalid URL",
			"operation", operation,
			"error", err)
		return entities.Targets{}, err
	}

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		s.logger.Errorw("Failed",
			"operation", operation,
			"error", err)
		return entities.Targets{}, err
	}

	dto, err := mapper.FromTarget(created)
	if err != nil {
		s.logger.Debugw("Failed",
			"operation", operation,
			"error", err)
	}

	reddisErr := s.cache.Set(ctx, created.URL, dto, DURATION)
	if reddisErr != nil {
		s.logger.Errorw("Failed",
			"operation", operation,
			"error", reddisErr.Error(),
		)
	}
	return created, nil
}

func (s *URLService) Get(ctx context.Context, id int) (entities.Targets, error) {
	key := strconv.Itoa(id)
	const operation = "get by id"
	data, err := s.cache.Get(ctx, key)
	//TODO: if cache hit
	if err == nil {
		dto, err := mapper.ToTarget(data)
		if err == nil {
			s.logger.Infoln("Succesful Redis",
				"operation", operation,
				"id", id)
			return dto, nil
		}
	}

	//singleflight через DoChan
	ch := s.group.DoChan(key, func() (interface{}, error) {
		data, err := s.cache.Get(ctx, key)
		if err == nil {
			dto, err := mapper.ToTarget(data)
			if err == nil {
				s.logger.Infoln("Succesful Redis",
					"operation", operation,
					"id", id)
				return dto, nil
			}
		}

		target, err := s.repo.GetByID(ctx, id)
		if err != nil {
			s.logger.Errorw("Failed",
				"operation", operation,
				"error", err)
			return entities.Targets{}, err
		}

		dto, err := mapper.FromTarget(target)
		if err == nil {
			_ = s.cache.Set(ctx, key, dto, DURATION)
		}
		s.logger.Debugw("Mapping failure",
			"operation", operation,
			"data", target)
		return target, nil
	})

	select {
	case res := <-ch:
		if res.Err != nil {
			s.logger.Errorw("Failed",
				"operation", operation,
				"error", err)
			return entities.Targets{}, res.Err
		}
		s.logger.Infow("Succes",
			"operation", operation)
		return res.Val.(entities.Targets), nil

	case <-ctx.Done():
		s.logger.Debugw("Channel was closed",
			"operation", operation)
		return entities.Targets{}, ctx.Err()
	}
}

func (s *URLService) List(ctx context.Context) (map[string]interface{}, error) {
	const operation = "get all"
	targets, err := s.cache.GetAll(ctx)
	if err == nil {
		s.logger.Errorw("Redis failure",
			"operation", operation,
			"error", err)
		return nil, err
	}
	return targets, nil
}

func (s *URLService) Update(ctx context.Context, id int, params entities.PatchReq) error {
	err := s.repo.Update(ctx, id, params)
	const operation = "update"
	if err != nil {
		s.logger.Errorw("Failed",
			"operation", operation,
			"error", err)
		return err
	}
	key := fmt.Sprintf("user:%v", id)
	_, err = s.cache.Delete(ctx, key)
	//TODO: if delete is failed
	if err == nil {
		s.logger.Debugw("Redis failure",
			"operation", operation,
			"error", err)
	}
	return nil
}

func (s *URLService) Delete(ctx context.Context, id int) error {
	_, err := s.cache.Delete(ctx, strconv.Itoa(id))
	const operation = "delete"
	if err == nil {
		s.logger.Debugw("Redis failure",
			"operation", operation,
			"error", err)
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed",
			"operation", operation,
			"error", err)
		return err
	}
	s.logger.Infoln("Succesful",
		"operation", operation)

	return nil
}

func (s *URLService) ListActive(ctx context.Context) ([]entities.Targets, error) {
	targets, err := s.repo.ListActive(ctx)

	const operation = "list active"
	if err != nil {
		s.logger.Errorw("Failed",
			"operation", operation,
			"error", err)

		return nil, err
	}
	s.logger.Infoln("Succesful",
		"operation", operation,
	)
	return targets, nil
}
