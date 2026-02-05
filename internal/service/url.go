package service

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/queries"
	"URL_checker/internal/service/validation"
	"context"
)

// Заменить название на TargetRepository
type IURLService interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id uint64) (entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
	Update(ctx context.Context, id uint64, params entities.PatchReq) error
	Delete(ctx context.Context, id uint64) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type URLService struct {
	repo *queries.PostgresRepo
}

func New(repo *queries.PostgresRepo) IURLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) Create(
	ctx context.Context,
	t entities.Targets,
) (entities.Targets, error) {

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return entities.Targets{}, err
	}

	return created, nil
}

func (s *URLService) Get(ctx context.Context, id uint64) (entities.Targets, error) {
	err := validation.ValidID(id)
	if err != nil {
		return entities.Targets{}, err
	}
	return s.repo.Get(ctx, id) //TODO: Посмотреть, на что ругается
}

func (s *URLService) List(ctx context.Context) ([]entities.Targets, error) {
	return s.repo.List(ctx) //TODO: Подумать, при каких условиях вывод изменится
}

func (s *URLService) Update(ctx context.Context, id uint64, params entities.PatchReq) error {
	err := validation.ValidID(uint64(id))
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, id, params)
}

func (s *URLService) Delete(ctx context.Context, id uint64) error {
	err := validation.ValidID(uint64(id))
	if err != nil {
		return err
	}
	s.repo.Delete(ctx, id)
	return nil
}

// TODO: после интеграции пострегесс вернуть все Targets с активным статусом, можно как-то сгруппировать
func (s *URLService) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return s.repo.ListActive(ctx)
}
