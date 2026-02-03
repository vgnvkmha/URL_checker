package service

import (
	entities "URL_checker/internal/repo/dto"
	"URL_checker/internal/service/validation"
)

// TODO: Подумать над возвращаемыми значениями
type URLService interface {
	PostTarget(url string, interval int, timeout int) bool
	GetTargets() []entities.URL
	GetTarget(id int) *entities.URL
	PatchTarget(params entities.PatchReq, id int) bool
	DeleteTarget(id int)
}

type urlService struct {
	activeUrls []entities.URL
}

func New() URLService {
	return &urlService{}
}

func (s *urlService) PostTarget(url string, interval int, timeout int) bool {
	if validation.Validation(url, interval, timeout) {
		newUrl := entities.URL{ID: len(s.activeUrls), URL: url, IntervalSec: interval, TimeoutMS: timeout}
		s.activeUrls = append(s.activeUrls, newUrl)
		return true
	}
	return false
}

func (s *urlService) GetTargets() []entities.URL {
	return s.activeUrls
}

func (s *urlService) GetTarget(id int) *entities.URL {
	if !validation.ValidID(id, len(s.activeUrls)) {
		return nil
	}
	return &s.activeUrls[id]
}

func (s *urlService) PatchTarget(params entities.PatchReq, id int) bool {
	target := s.GetTarget(id)
	if target == nil {
		return false
	}
	if params.Interval != nil {
		if !validation.IsValidInterval(*params.Interval) {
			return false
		}
		target.IntervalSec = *params.Interval
	}

	if params.Timeout != nil {
		if !validation.IsValidTimeout(*params.Timeout) {
			return false
		}
		target.TimeoutMS = *params.Timeout
	}

	if params.Active != nil {
		target.Active = *params.Active
	}

	return true
}

func (s *urlService) DeleteTarget(id int) {
	if validation.ValidID(id, len(s.activeUrls)) {
		s.activeUrls = append(s.activeUrls[:id], s.activeUrls[id+1:]...)
	}
}
