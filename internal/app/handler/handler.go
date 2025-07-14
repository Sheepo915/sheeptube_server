package handler

import (
	"sheeptube/internal/app/service"
)

type Handler struct {
	service *service.Service
	VideoHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service:      service,
		VideoHandler: newVideoHandler(service.VideoService),
	}
}
