package api

import (
	"encoding/json"
	"io"
	"mobile-connect/internal/apk/entity"
	"net/http"
	"os"
	"strings"
)

func (h *Handler) InstallAPK(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/devices/")
	id = strings.TrimSuffix(id, "/install")

	file, header, err := r.FormFile("apk")
	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp(
		"",
		"*.apk",
	)
	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(
		tmpFile,
		file,
	)
	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	output, err := h.apkService.InstallAPK(
		id,
		tmpFile.Name(),
	)

	response := map[string]string{
		"file":   header.Filename,
		"output": output,
	}

	if err != nil {

		w.WriteHeader(
			http.StatusInternalServerError,
		)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UninstallAPK(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/devices/")
	id = strings.TrimSuffix(id, "/uninstall")

	var req entity.UninstallRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	output, err := h.apkService.UninstallAPK(
		id,
		req.Package,
	)

	response := map[string]string{
		"package": req.Package,
		"output":  output,
	}

	if err != nil {
		w.WriteHeader(
			http.StatusInternalServerError,
		)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetPackages(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/devices/")
	id = strings.TrimSuffix(id, "/packages")

	packages, err := h.apkService.ListPackages(id)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		packages,
	)
}

func (h *Handler) LaunchApp(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	var req entity.AppRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	err := h.apkService.LaunchApp(
		id,
		req.Package,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"device":  id,
			"package": req.Package,
			"status":  "launched",
		},
	)
}

func (h *Handler) StopApp(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	var req entity.AppRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	err := h.apkService.StopApp(
		id,
		req.Package,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"device":  id,
			"package": req.Package,
			"status":  "stopped",
		},
	)
}
