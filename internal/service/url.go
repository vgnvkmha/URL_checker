package service

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/queries"
	"URL_checker/internal/service/validation"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Заменить название на TargetRepository
type URLService interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id uint64) (entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
	Update(ctx context.Context, id uint64, params entities.PatchReq) error
	Delete(ctx context.Context, id uint64) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type urlService struct {
	repo       *queries.PostgresRepo
	activeUrls []entities.Targets
}

func New(repo *queries.PostgresRepo) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) Create(
	ctx context.Context,
	t entities.Targets,
) (entities.Targets, error) {

	select {
	case <-ctx.Done():
		return entities.Targets{}, ctx.Err()
	default:
	}

	err := validation.Validation(t.URL, t.IntervalSec, t.TimeoutMS)
	if err != nil {
		return entities.Targets{}, err
	}

	s.repo.Create(ctx, t)

	return t, nil //TODO: Посмотреть, на что ругается
}

func (s *urlService) Get(ctx context.Context, id uint64) (entities.Targets, error) {
	err := validation.ValidID(id, len(s.activeUrls))
	if err != nil {
		return entities.Targets{}, err
	}
	return s.activeUrls[id], nil //TODO: Посмотреть, на что ругается
}

func (s *urlService) List(ctx context.Context) ([]entities.Targets, error) {
	return s.activeUrls, nil //TODO: Подумать, при каких условиях вывод изменится
}

func (s *urlService) Update(ctx context.Context, id uint64, params entities.PatchReq) error {
	target, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if params.Interval != nil {
		err := validation.IsValidInterval(*params.Interval)
		if err != nil {
			return err
		}
		target.IntervalSec = *params.Interval
		target.UpdatedAt = *timestamppb.Now()
	}

	if params.Timeout != nil {
		err := validation.IsValidTimeout(*params.Interval)
		if err != nil {
			return err
		}
		target.TimeoutMS = *params.Timeout
		target.UpdatedAt = *timestamppb.Now()
	}

	if params.Active != nil {
		target.Active = *params.Active
		target.UpdatedAt = *timestamppb.Now()
	}

	return nil
}

func (s *urlService) Delete(ctx context.Context, id uint64) error {
	err := validation.ValidID(uint64(id), len(s.activeUrls))
	if err != nil {
		return err
	}
	s.activeUrls = append(s.activeUrls[:id], s.activeUrls[id+1:]...)
	return nil
}

// TODO: после интеграции пострегесс вернуть все Targets с активным статусом, можно как-то сгруппировать
func (s *urlService) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return nil, nil
}
