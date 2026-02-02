package service

import (
	"URL_checker/internal/repo/entities"
	"URL_checker/internal/service/validation"
)

// TODO: Подумать над возвращаемыми значениями
type URLService interface {
	postTarget(url string, interval int, timeout int)
	getTargets() []entities.URL
	getTarget(id int) entities.URL
	patchTarget(url string, interval int, timeout int, id int)
	deleteTarget(id int)
}

type urlService struct {
	activeUrls []entities.URL
}

func (s *urlService) postTarget(url string, interval int, timeout int) {
	if validation.Validation(url, interval, timeout) {
		newUrl := entities.URL{URL: url, IntervalSec: interval, TimeoutMS: timeout}
		s.activeUrls = append(s.activeUrls, newUrl)
	}
}

func (s *urlService) getTargets() []entities.URL {
	return s.activeUrls
}

func (s *urlService) getTarget(id int) entities.URL {
	if validation.ValidID(id, len(s.activeUrls)) {
		return s.activeUrls[id-1]
	} else {
		//TODO: Убрать заглушку
		return s.activeUrls[0]
	}
}

// TODO: Добавить возможность получать не все параметры
func (s *urlService) patchTarget(url string, interval int, timeout int, id int) {
	if validation.Validation(url, interval, timeout) && validation.ValidID(id, len(s.activeUrls)) {
		s.activeUrls[id].IntervalSec = interval
		s.activeUrls[id].URL = url
		s.activeUrls[id].TimeoutMS = timeout
	}
}

func (s *urlService) deleteTarget(id int) {
	if validation.ValidID(id, len(s.activeUrls)) {
		s.activeUrls = append(s.activeUrls[:id], s.activeUrls[id+1:]...)
	}
}
