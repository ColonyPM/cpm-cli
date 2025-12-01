package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type manifestSchema struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
}

func validateManifest(manifestPath string) (manifestSchema, error) {
	var manifest manifestSchema

	info, err := os.Stat(manifestPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return manifest, fmt.Errorf("%s not found; run `cpm init` first", manifestName)
		}
		return manifest, fmt.Errorf("failed to get manifest file: %v", err)
	}
	if !info.Mode().IsRegular() {
		return manifest, fmt.Errorf("%s is not a regular file", manifestName)
	}

	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return manifest, fmt.Errorf("failed to read %s: %v", manifestName, err)
	}

	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return manifest, fmt.Errorf("failed to parse %s: %v", manifestName, err)
	}

	var missing []string
	if manifest.Name == "" {
		missing = append(missing, "name")
	}
	if manifest.Version == "" {
		missing = append(missing, "version")
	}
	if manifest.Description == "" {
		missing = append(missing, "description")
	}
	if manifest.Author == "" {
		missing = append(missing, "author")
	}

	if len(missing) > 0 {
		return manifest, fmt.Errorf("%s missing required field(s): %s", manifestName, strings.Join(missing, ", "))
	}

	return manifest, nil
}

func addDirToArchive(dir string, root string, tw *tar.Writer) error {
	base := filepath.Clean(dir)
	rootName := filepath.Clean(root)
	if rootName == "." || rootName == "" {
		rootName = filepath.Base(base)
	}

	return filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		hdr.Name = filepath.ToSlash(filepath.Join(rootName, rel))
		if info.IsDir() && hdr.Name[len(hdr.Name)-1] != '/' {
			hdr.Name += "/"
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
}

func archiveDirectory(dir string, root string) ([]byte, error) {
	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// write all files into the tar
	if err := addDirToArchive(dir, root, tw); err != nil {
		tw.Close()
		gw.Close()
		return nil, err
	}

	// close writers to flush everything into buf
	if err := tw.Close(); err != nil {
		gw.Close()
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func uploadPackage(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get cwd: %v", err)
	}

	manifestPath := filepath.Join(cwd, manifestName)
	manifest, err := validateManifest(manifestPath)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	archive, err := archiveDirectory(cwd, manifest.Name)
	if err != nil {
		return err
	}

	url, err := repoClient.Upload(ctx, manifest.Name, archive, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Upload successful: %s\n", url)

	return nil
}

var uploadCmd = &cobra.Command{
	Use:   "upload <upload-token>",
	Short: "Upload a package to the repository",
	Long:  "Upload a package to the repository",
	Args:  cobra.ExactArgs(1),
	RunE:  uploadPackage,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
