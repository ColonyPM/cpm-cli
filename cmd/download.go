/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	downloadVersion  string
	downloadOutput   string
	downloadRegistry string
	downloadExtract  bool
)

func downloadPackage(packageName, version, outputDir, registry string, extract bool) error {
	pkgName := args[0]
	if pkgName == "" {
		return errors.New("package name is required")
	}

	// Resolve default registry url from config/env later if you want.
	regURL := strings.TrimRight(downloadRegistry, "/")
	if regURL == "" {
		return errors.New("registry URL is empty")
	}

	version := downloadVersion
	if version == "" || version == "latest" {
		v, err := resolveLatestVersion(regURL, pkgName)
		if err != nil {
			return fmt.Errorf("failed to resolve latest version: %w", err)
		}
		version = v
	}

	if err := os.MkdirAll(downloadOutput, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	archivePath := filepath.Join(downloadOutput, fmt.Sprintf("%s-%s.cpm.tgz", pkgName, version))

	if err := downloadPackage(regURL, pkgName, version, archivePath); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s@%s to %s\n", pkgName, version, archivePath)

	if downloadExtract {
		if err := extractTarGz(archivePath, downloadOutput); err != nil {
			return fmt.Errorf("failed to extract archive: %w", err)
		}
		fmt.Printf("Extracted to %s\n", downloadOutput)
	}
	return nil
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <package>",
	Short: "Download a package from the registry",
	Long: `Download a ColonyOS package from the registry.

By default this will download the specified package archive.
You can specify a version and output directory, and choose
whether to auto-extract the archive.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		return downloadPackage(args[0], downloadVersion, downloadOutput, downloadRegistry, downloadExtract)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&downloadVersion, "version", "v", "latest", "Package version (default: latest)")
	downloadCmd.Flags().StringVarP(&downloadOutput, "output", "o", ".", "Output directory")
	downloadCmd.Flags().StringVar(&downloadRegistry, "registry", "https://registry.colonyos.local", "Registry URL")
	downloadCmd.Flags().BoolVarP(&downloadExtract, "extract", "x", false, "Extract the downloaded archive")
}
