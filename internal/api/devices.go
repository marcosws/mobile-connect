package api

import (
	"encoding/json"
	"mobile-connect/internal/logger"
	"net/http"
	"strings"
)

func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {

	devices, err := h.deviceService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/devices/")

	logger.Logger.Printf("searching device id=%s", id)

	device, err := h.deviceService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(device)
}

// func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
// 	id := strings.TrimPrefix(
// 		r.URL.Path,
// 		"/devices/",
// 	)
// 	logger.Logger.Printf(
// 		"searching device id=%s",
// 		id,
// 	)
// 	fmt.Fprintf(w, "device: %s", id)
// }
