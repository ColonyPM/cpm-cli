package cmd

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func getDestPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return "./pkgs", nil // %LOCALAPPDATA%\cpm\pkgs
	case "linux", "freebsd", "openbsd", "netbsd":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("locate home directory: %w", err)
		}
		return filepath.Join(home, ".local", "share", "cpm", "pkgs"), nil // $XDG_DATA_HOME
	case "darwin":
		return "", fmt.Errorf("macOS support not implemented 💀") // TODO: Implement macOS support (~/Library/Application Support/YourApp?)
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func extractArchive(rr io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(rr)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read tar entry: %w", err)
		}

		targetPath := filepath.Join(destDir, hdr.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", targetPath)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(hdr.Mode)); err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return fmt.Errorf("mkdir parents: %w", err)
			}

			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("write file: %w", err)
			}
			f.Close()
		}
	}
}

func downloadPackage(cmd *cobra.Command, args []string) error {
	fmt.Println(" Downloading " + args[0] + "...")
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	destDir, err := getDestPath()
	if err != nil {
		return err
	}
	destDir, err = filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("resolve destination path: %w", err)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("create destination path: %w", err)
	}

	rr, err := repoClient.Download(ctx, args[0])
	if err != nil {
		return err
	}
	defer rr.Close()

	if err := extractArchive(rr, destDir); err != nil {
		return err
	}

	fmt.Println(" Downloaded " + args[0] + " ✅")
	return nil
}

var downloadCmd = &cobra.Command{
	Use:   "download <NAME>",
	Short: "Download a package from the repository",
	Long:  "Download a package from the repository",
	Args:  cobra.ExactArgs(1),
	RunE:  downloadPackage,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
