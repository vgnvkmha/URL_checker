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
	repo  target.ITargetRepo
	cache cache.IRedisCache
}

func New(repo target.ITargetRepo, cache cache.IRedisCache) *URLService {
	return &URLService{
		repo:  repo,
		cache: cache,
	}
}

func (s *URLService) Create(
	ctx context.Context,
	t entities.Targets,
) (entities.Targets, error) {

	err := validation.ValidURL(t.URL)
	if err != nil {
		return entities.Targets{}, err
	}

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return entities.Targets{}, err
	}

	dto, _ := mapper.FromTarget(created)

	reddisErr := s.cache.Set(ctx, created.URL, dto, DURATION)
	if reddisErr != nil {
		fmt.Println("SET REDIS ERROR:", reddisErr.Error())
	}

	return created, nil
}

func (s *URLService) Get(ctx context.Context, id int) (entities.Targets, error) {
	key := strconv.Itoa(id)

	data, err := s.cache.Get(ctx, key)
	if err == nil {
		dto, err := mapper.ToTarget(data)
		if err == nil {
			return dto, nil
		}
		// если ошибка десериализации — идём в БД
	}

	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entities.Targets{}, err
	}

	dto, _ := mapper.FromTarget(target)
	err = s.cache.Set(ctx, key, dto, DURATION)

	return target, err
}

func (s *URLService) List(ctx context.Context) (map[string]interface{}, error) {
	return s.cache.GetAll(ctx)
}

func (s *URLService) Update(ctx context.Context, id int, params entities.PatchReq) error {
	err := s.repo.Update(ctx, id, params)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("user:%v", id)
	_, err = s.cache.Delete(ctx, key)
	if err != nil {
		fmt.Println("REDIS DELETE SUCCESFUL")
	}
	return nil
}

func (s *URLService) Delete(ctx context.Context, id int) error {
	status, err := s.cache.Delete(ctx, strconv.Itoa(id))
	if err == nil {
		fmt.Println("REDIS DELETE ERROR", status)
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLService) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return s.repo.ListActive(ctx)
}
