package url

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
	"context"
)

type IURLHandler interface {
	Create(ctx context.Context, t entities.Targets) (entities.Targets, error)
	Get(ctx context.Context, id uint64) (entities.Targets, error)
	List(ctx context.Context) ([]entities.Targets, error)
	Update(ctx context.Context, id uint64, params entities.PatchReq) error
	Delete(ctx context.Context, id uint64) error
	ListActive(ctx context.Context) ([]entities.Targets, error)
}

type URLHandler struct {
	service service.IURLService
}

func New(service service.IURLService) IURLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) Create(ctx context.Context, t entities.Targets) (entities.Targets, error) {
	return h.service.Create(ctx, t)
}

func (h *URLHandler) Get(ctx context.Context, id uint64) (entities.Targets, error) {
	return h.service.Get(ctx, id)
}

func (h *URLHandler) List(ctx context.Context) ([]entities.Targets, error) {
	return h.service.List(ctx)
}

func (h *URLHandler) Update(ctx context.Context, id uint64, params entities.PatchReq) error {
	return h.service.Update(ctx, id, params)
}

func (h *URLHandler) ListActive(ctx context.Context) ([]entities.Targets, error) {
	return h.service.ListActive(ctx)
}

func (h *URLHandler) Delete(ctx context.Context, id uint64) error {
	return h.service.Delete(ctx, id)
}
