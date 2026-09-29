package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mobile-connect/internal/appium"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddAppiumUDID(t *testing.T) {
	body, err := addAppiumUDID([]byte(`{"capabilities":{"alwaysMatch":{"platformName":"Android"}},"desiredCapabilities":{"udid":"other"}}`), "emulator-5554")
	if err != nil {
		t.Fatalf("add udid: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	w3c := got["capabilities"].(map[string]any)["alwaysMatch"].(map[string]any)
	if w3c["appium:udid"] != "emulator-5554" {
		t.Errorf("W3C udid = %v, want emulator-5554", w3c["appium:udid"])
	}
	legacy := got["desiredCapabilities"].(map[string]any)
	if legacy["udid"] != "emulator-5554" {
		t.Errorf("legacy udid = %v, want emulator-5554", legacy["udid"])
	}
}

func TestCreateAppiumSessionForwardsCapabilities(t *testing.T) {
	var receivedPath string
	var received map[string]any
	client, err := appium.NewWithHTTPClient("http://appium.test", &http.Client{Transport: roundTripper(func(r *http.Request) (*http.Response, error) {
		receivedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode upstream request: %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"value":{"sessionId":"abc"}}`)),
		}, nil
	})})
	if err != nil {
		t.Fatalf("new Appium client: %v", err)
	}
	handler := &Handler{appiumClient: client}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /devices/{id}/appium/sessions", handler.CreateAppiumSession)

	request := httptest.NewRequest(http.MethodPost, "/devices/device-42/appium/sessions", strings.NewReader(`{"capabilities":{"alwaysMatch":{"platformName":"Android"}}}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		body, _ := io.ReadAll(response.Result().Body)
		t.Fatalf("status = %d, want 200; body = %s", response.Code, body)
	}
	if receivedPath != "/session" {
		t.Errorf("upstream path = %q, want /session", receivedPath)
	}
	alwaysMatch := received["capabilities"].(map[string]any)["alwaysMatch"].(map[string]any)
	if alwaysMatch["appium:udid"] != "device-42" {
		t.Errorf("udid = %v, want device-42", alwaysMatch["appium:udid"])
	}
}

func TestAppiumProxyRejectsDirectSessionCreation(t *testing.T) {
	handler := &Handler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/appium/{path...}", handler.AppiumProxy)

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/appium/session", nil))

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
