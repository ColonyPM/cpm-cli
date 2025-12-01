package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Upload(ctx context.Context, path, filePath string) error {
	return nil
}

func (c *Client) Download(ctx context.Context, packageName string) (io.ReadCloser, error) {
	path := fmt.Sprintf("packages/%s/download", packageName)
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10)) // 4KB
		return nil, fmt.Errorf("download failed: status=%d body=%q", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}
