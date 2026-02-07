package service

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"context"
)

type IURLService interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id int) (entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
	Update(ctx context.Context, id int, params entities.PatchReq) error
	Delete(ctx context.Context, id int) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type URLService struct {
	repo target.ITargetRepo
}

func New(repo target.ITargetRepo) *URLService {
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

func (s *URLService) Get(ctx context.Context, id int) (entities.Targets, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *URLService) List(ctx context.Context) ([]entities.Targets, error) {
	return s.repo.List(ctx)
}

func (s *URLService) Update(ctx context.Context, id int, params entities.PatchReq) error {
	return s.repo.Update(ctx, id, params)
}

func (s *URLService) Delete(ctx context.Context, id int) error {
	s.repo.Delete(ctx, id)
	return nil
}

func (s *URLService) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return s.repo.ListActive(ctx)
}
