package api

import (
	"encoding/json"
	"fmt"
	"mobile-connect/internal/devices"
	"mobile-connect/internal/logger"
	"net/http"
	"strings"
)

type Handler struct {
	deviceService *devices.Service
}

func NewHandler(
	deviceService *devices.Service,
) *Handler {

	return &Handler{
		deviceService: deviceService,
	}
}

func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {

	devices, err := h.deviceService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) GetDeviceByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(
		r.URL.Path,
		"/devices/",
	)
	logger.Logger.Printf(
		"searching device id=%s",
		id,
	)
	fmt.Fprintf(w, "device: %s", id)
}
