package repository

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (c *Client) Download(ctx context.Context, packageName, destPath string) error {
	fmt.Println("Downloading package...")

	// Create URL
	path := fmt.Sprintf("packages/%s/download", packageName)
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	// Send the request
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Fail if status is not 2xx.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read a small part of the body for a helpful error message.
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10)) // 4KB
		return fmt.Errorf("download failed: status=%d body=%q", resp.StatusCode, string(body))
	}
	fmt.Println("Got HTTP status:", resp.StatusCode)

	// Make sure the destination directory exists.
	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return fmt.Errorf("create dest dir: %w", err)
	}

	// Create a temporary file to store the .tar.gz
	tmpFile, err := os.CreateTemp("", "cpm-*.tar.gz")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // clean up after extraction
	}()

	// Stream response into the temp file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("save archive: %w", err)
	}

	// Rewind to the beginning so it can be read again
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind archive: %w", err)
	}

	fmt.Println("Archive downloaded to:", tmpFile.Name())

	// Create a gzip reader
	gzReader, err := gzip.NewReader(tmpFile)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create a tar reader on top of the gzip reader
	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// end of archive
			break
		}
		if err != nil {
			return fmt.Errorf("read tar entry: %w", err)
		}

		// Clean the name to avoid weird paths like ../../etc/passwd
		targetPath := filepath.Join(destPath, filepath.Clean(header.Name))

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("create dir %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// Ensure containing directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return fmt.Errorf("create parent dir %s: %w", filepath.Dir(targetPath), err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("write file %s: %w", targetPath, err)
			}
			outFile.Close()

		default:
			// You can add support for other types (symlinks, etc.) later
			fmt.Println("Skipping unsupported entry:", header.Name)
		}
	}

	fmt.Println("Extracted archive into:", destPath)
	return nil
}
