package api

import (
	"mobile-connect/internal/devices"
	"mobile-connect/internal/shell"
)

type Handler struct {
	deviceService *devices.Service
	shellService  *shell.Service
}

func NewHandler(
	deviceService *devices.Service,
	shellService *shell.Service,
) *Handler {

	return &Handler{
		deviceService: deviceService,
		shellService:  shellService,
	}
}
