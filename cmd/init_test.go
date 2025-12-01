package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitPackageCreatesManifest(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	t.Setenv("USER", "testuser")

	if err := initPackage(&cobra.Command{}, nil); err != nil {
		t.Fatalf("initPackage() error = %v", err)
	}

	path := filepath.Join(tmp, manifestName)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected manifest to be created: %v", err)
	}

	got := string(data)
	if !strings.Contains(got, `name: "`+filepath.Base(tmp)+`"`) {
		t.Errorf("manifest missing package name, got:\n%s", got)
	}
}

func TestInitManifestAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	_, err := os.Create(manifestName)
	if err != nil {
		t.Fatalf("failed to create manifest: %v", err)
	}

	err = initPackage(&cobra.Command{}, nil)
	if err == nil {
		t.Fatalf("initPackage() error = nil, want ErrManifestAlreadyExists")
	}

	if !errors.Is(err, ErrManifestAlreadyExists) {
		t.Fatalf("initPackage() error = %v, want %v", err, ErrManifestAlreadyExists)
	}
}
