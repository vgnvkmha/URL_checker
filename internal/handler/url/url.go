package url

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
	"context"
)

type URLHandler interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id uint64) (entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
	Update(ctx context.Context, id uint64, params entities.PatchReq) error
	Delete(ctx context.Context, id uint64) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type URLController struct {
	service service.URLService
}

func New(service service.URLService) URLHandler {
	return &URLController{
		service: service,
	}
}

func (h *URLController) Create(ctx context.Context, t entities.Targets) (entities.Targets, error) {
	return h.service.Create(ctx, t)
}

func (h *URLController) Get(ctx context.Context, id uint64) (entities.Targets, error) {
	return h.service.Get(ctx, id)
}

func (h *URLController) List(ctx context.Context) ([]entities.Targets, error) {
	return h.service.List(ctx)
}

func (h *URLController) Update(ctx context.Context, id uint64, params entities.PatchReq) error {
	return h.service.Update(ctx, id, params)
}

func (h *URLController) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return h.service.ListActive(ctx)
}

func (h *URLController) Delete(ctx context.Context, id uint64) error {
	return h.service.Delete(ctx, id)
}
