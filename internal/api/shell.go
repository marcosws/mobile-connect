package api

import (
	"encoding/json"
	"mobile-connect/internal/devices/entity"
	"net/http"
	"strings"
)

func (h *Handler) Shell(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := strings.TrimPrefix(r.URL.Path, "/devices/")
	id = strings.TrimSuffix(id, "/shell")

	var req entity.ShellRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	output, err := h.shellService.Shell(
		id,
		req.Command,
	)

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
		map[string]string{
			"output": output,
		},
	)
}
