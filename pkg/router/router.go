package router

import (
	"download_bot/pkg/interfaces"
	"errors"
)

type Router struct {
	services []interfaces.VideoService
}

func NewRouter(services ...interfaces.VideoService) *Router {
	return &Router{services: services}
}

func (r *Router) Resolve(url string) (interfaces.VideoService, error) {
	for _, s := range r.services {
		if s.Match(url) == true {
			return s, nil
		}
	}
	return nil, errors.New("Неизвестный сервис")
}
