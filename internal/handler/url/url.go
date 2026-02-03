package url

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service"
)

type URLHandler interface {
	PostTarget(url string, interval int, timeout int) bool
	GetTargets() []entities.URL
	GetTarget(id int) *entities.URL
	PatchTarget(params entities.PatchReq, id int) bool
	DeleteTarget(id int)
}

type URLController struct {
	service service.URLService
}

func New(service service.URLService) URLHandler {
	return &URLController{
		service: service,
	}
}

func (h *URLController) PostTarget(url string, interval int, timeout int) bool {
	return h.service.PostTarget(url, interval, timeout)
}

func (h *URLController) GetTargets() []entities.URL {
	return h.service.GetTargets()
}

func (h *URLController) GetTarget(id int) *entities.URL {
	return h.service.GetTarget(id)
}

func (h *URLController) PatchTarget(params entities.PatchReq, id int) bool {
	if h.service.PatchTarget(params, id) {
		return true
	}
	return false
}

func (h *URLController) DeleteTarget(id int) {
	h.service.DeleteTarget(id)
}
