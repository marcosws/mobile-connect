package api

import (
	"mobile-connect/internal/devices"
	"mobile-connect/internal/shell"
	"mobile-connect/internal/stream"
)

type Handler struct {
	deviceService *devices.Service
	shellService  *shell.Service
	streamService *stream.Service
}

func NewHandler(
	deviceService *devices.Service,
	shellService *shell.Service,
	streamService *stream.Service,
) *Handler {

	return &Handler{
		deviceService: deviceService,
		shellService:  shellService,
		streamService: streamService,
	}
}
