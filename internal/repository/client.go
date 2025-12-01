package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

func (c *Client) Upload(ctx context.Context, pkgName string, data []byte, uploadToken string) (string, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, "packages/upload")

	// Build multipart/form-data body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// field name MUST be "archive" to match `archive: UploadedFile = File(...)`
	fw, err := writer.CreateFormFile("archive", pkgName+".tar.gz")
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := fw.Write(data); err != nil {
		return "", fmt.Errorf("write archive data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", uploadToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var result struct {
			Detail string `json:"detail"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", fmt.Errorf("decode upload response: %w", err)
		}
		return "", fmt.Errorf("upload failed: %s", result.Detail)
	}

	var ok struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ok); err != nil {
		return "", fmt.Errorf("decode upload response: %w", err)
	}

	return ok.URL, nil
}
