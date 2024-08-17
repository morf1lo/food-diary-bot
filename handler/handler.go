package handler

import (
	"time"

	"github.com/morf1lo/food-diary-bot/service"
)

const cooldownDuration = time.Minute

type Handler struct {
	services *service.Service
	lastCommands map[int64]map[string]time.Time
}

func New(services *service.Service) *Handler {
	return &Handler{
		services: services,
		lastCommands: make(map[int64]map[string]time.Time),
	}
}
