package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"omniport-api/internal/middleware"
)

type InternalRESTClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewInternalRESTClient(baseURL string, token string, timeout time.Duration) *InternalRESTClient {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &InternalRESTClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		token:   strings.TrimSpace(token),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *InternalRESTClient) Enabled() bool {
	return c != nil && c.baseURL != "" && c.client != nil
}

func (c *InternalRESTClient) PostJSON(ctx context.Context, path string, payload interface{}, headers map[string]string) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set(middleware.InternalServiceTokenHeader, c.token)
	}
	for key, value := range headers {
		if strings.TrimSpace(value) != "" {
			req.Header.Set(key, value)
		}
	}

	return c.client.Do(req)
}
