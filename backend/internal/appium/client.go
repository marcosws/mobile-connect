// Package appium contains the HTTP client used to communicate with an Appium server.
package appium

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Client forwards WebDriver requests to a single, configured Appium server.
// BaseURL may include a path, for example http://localhost:4723/wd/hub.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func New(baseURL string) (*Client, error) {
	return NewWithHTTPClient(baseURL, http.DefaultClient)
}

// NewWithHTTPClient creates a client with a custom HTTP client. It is useful when
// callers need to configure timeouts, authentication transports, or tests.
func NewWithHTTPClient(baseURL string, httpClient *http.Client) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Appium URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid Appium URL: use an http(s) URL with a host")
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: parsed, httpClient: httpClient}, nil
}

// Forward sends a request to Appium and returns its original response. requestPath
// must be relative to Appium's configured base path.
func (c *Client) Forward(ctx context.Context, method, requestPath string, body []byte, headers http.Header) (*http.Response, error) {
	if !strings.HasPrefix(requestPath, "/") {
		requestPath = "/" + requestPath
	}

	target := *c.baseURL
	target.Path = path.Join(c.baseURL.Path, requestPath)
	if strings.HasSuffix(requestPath, "/") && !strings.HasSuffix(target.Path, "/") {
		target.Path += "/"
	}

	req, err := http.NewRequestWithContext(ctx, method, target.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Appium request: %w", err)
	}
	for key, values := range headers {
		if strings.EqualFold(key, "Host") || strings.EqualFold(key, "Content-Length") {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connect to Appium: %w", err)
	}
	return response, nil
}

// BaseURL returns a copy of the configured Appium endpoint.
func (c *Client) BaseURL() url.URL {
	return *c.baseURL
}

// Available reports whether the configured Appium server responds to its status
// endpoint with a successful HTTP status.
func (c *Client) Available(ctx context.Context) bool {
	response, err := c.Forward(ctx, http.MethodGet, "/status", nil, nil)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
}

// CopyResponse copies status, headers and body without changing Appium's WebDriver
// response format.
func CopyResponse(w http.ResponseWriter, response *http.Response) error {
	defer response.Body.Close()
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(response.StatusCode)
	_, err := io.Copy(w, response.Body)
	return err
}
