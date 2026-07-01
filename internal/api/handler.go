package api

import (
	"mobile-connect/internal/apk"
	"mobile-connect/internal/apps"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/shell"
	"mobile-connect/internal/stream"
)

type Handler struct {
	deviceService *devices.Service
	shellService  *shell.Service
	streamService *stream.Service
	apkService    *apk.Service
	appService    *apps.Service
}

func NewHandler(
	deviceService *devices.Service,
	shellService *shell.Service,
	streamService *stream.Service,
	apkService *apk.Service,
	appService *apps.Service,
) *Handler {

	return &Handler{
		deviceService: deviceService,
		shellService:  shellService,
		streamService: streamService,
		apkService:    apkService,
		appService:    appService,
	}
}
