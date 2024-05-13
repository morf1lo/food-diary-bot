package handler

import "github.com/morf1lo/food-diary-bot/service"

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}
