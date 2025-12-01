package repository

import (
	"context"
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

func (c *Client) Download(ctx context.Context, path, destPath string) error {
	return nil
}
