package group

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client fetches group data from the upstream group service.
type Client interface {
	ListGroups(ctx context.Context) ([]Group, error)
}

type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient builds a new API client using the provided base URL and timeout.
func NewHTTPClient(baseURL string, timeout time.Duration) Client {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	trimmed := strings.TrimRight(baseURL, "/")
	return &httpClient{
		baseURL: trimmed,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *httpClient) ListGroups(ctx context.Context) ([]Group, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/group/", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request group service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("group service returned status %d", resp.StatusCode)
	}

	var payload struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Data    []Group `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode group response: %w", err)
	}

	if !payload.Success {
		return nil, fmt.Errorf("group service error: %s", strings.TrimSpace(payload.Message))
	}

	return payload.Data, nil
}
