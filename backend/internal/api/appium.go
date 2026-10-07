package api

import (
	"encoding/json"
	"errors"
	"io"
	"mobile-connect/internal/appium"
	"net/http"
	"strings"
)

const maxAppiumRequestBody = 8 << 20

// AppiumStatus proxies Appium's status endpoint. It is useful for checking whether
// the automation server is ready before creating a session.
func (h *Handler) AppiumStatus(w http.ResponseWriter, r *http.Request) {
	h.forwardToAppium(w, r, "/status", nil)
}

// CreateAppiumSession creates a W3C WebDriver session and forces the selected
// Mobile Connect device as Appium's UDID. This prevents a client from accidentally
// running a request against a different connected phone.
func (h *Handler) CreateAppiumSession(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	if deviceID == "" {
		http.Error(w, "device id is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxAppiumRequestBody))
	if err != nil {
		http.Error(w, "invalid Appium request: "+err.Error(), http.StatusBadRequest)
		return
	}

	body, err = addAppiumUDID(body, deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.forwardToAppium(w, r, "/session", body)
}

// AppiumProxy forwards session commands (find element, click, execute script, and
// so on) to Appium. The route intentionally exposes only relative WebDriver paths.
func (h *Handler) AppiumProxy(w http.ResponseWriter, r *http.Request) {
	h.proxyAppiumPath(w, r, r.PathValue("path"))
}

// AppiumDeviceProxy is the WebDriver-compatible route used by clients such as
// Appium Inspector. Its base path includes the device ID selected for the session.
func (h *Handler) AppiumDeviceProxy(w http.ResponseWriter, r *http.Request) {
	h.proxyAppiumPath(w, r, r.PathValue("path"))
}

func (h *Handler) proxyAppiumPath(w http.ResponseWriter, r *http.Request, pathValue string) {
	requestPath := "/" + strings.TrimPrefix(pathValue, "/")
	if requestPath == "/" {
		http.NotFound(w, r)
		return
	}
	if requestPath == "/session" && r.Method == http.MethodPost {
		http.Error(w, "create sessions through POST /devices/{id}/appium/sessions", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxAppiumRequestBody))
	if err != nil {
		http.Error(w, "invalid Appium request: "+err.Error(), http.StatusBadRequest)
		return
	}
	h.forwardToAppium(w, r, requestPath, body)
}

func (h *Handler) forwardToAppium(w http.ResponseWriter, r *http.Request, requestPath string, body []byte) {
	response, err := h.appiumClient.Forward(r.Context(), r.Method, requestPath, body, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := appium.CopyResponse(w, response); err != nil {
		return
	}
}

func addAppiumUDID(body []byte, deviceID string) ([]byte, error) {
	var capabilities map[string]any
	if err := json.Unmarshal(body, &capabilities); err != nil {
		return nil, errors.New("body must be a valid Appium capabilities JSON object")
	}

	if capabilities == nil {
		return nil, errors.New("body must be an Appium capabilities JSON object")
	}

	// W3C capability format used by Appium 2.
	w3c, ok := capabilities["capabilities"].(map[string]any)
	if !ok {
		w3c = make(map[string]any)
		capabilities["capabilities"] = w3c
	}
	alwaysMatch, ok := w3c["alwaysMatch"].(map[string]any)
	if !ok {
		alwaysMatch = make(map[string]any)
		w3c["alwaysMatch"] = alwaysMatch
	}
	alwaysMatch["appium:udid"] = deviceID

	// Keeping this also supports Appium 1 clients that still send JSONWP desiredCapabilities.
	legacy, ok := capabilities["desiredCapabilities"].(map[string]any)
	if ok {
		legacy["udid"] = deviceID
	}

	encoded, err := json.Marshal(capabilities)
	if err != nil {
		return nil, errors.New("could not encode Appium capabilities")
	}
	return encoded, nil
}
