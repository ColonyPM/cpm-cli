package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const manifestTemplate = `name: "{{ .Name }}"
version: "0.0.0"
description: ""
author: ""
`

type ManifestData struct {
	Name string
}

func createManifestFile(dir string, name string) error {
	path := filepath.Join(dir, "package.yml")

	tmpl, err := template.New("manifest").Parse(manifestTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	data := ManifestData{Name: name}
	if err := tmpl.Execute(f, data); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new package",
	Long:  "cpm init - Initialize a new package in the cwd \n\n ",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			if err = createManifestFile(dir, filepath.Base(dir)); err != nil {
				return fmt.Errorf("create manifest: %w", err)
			}
		} else {
			dirName := args[0]

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			dir := filepath.Join(wd, dirName)

			if err := os.Mkdir(dir, 0o755); err != nil {
				if os.IsExist(err) {
					return fmt.Errorf("directory %q already exists", dir)
				}
				return err
			}

			if err := createManifestFile(dir, dirName); err != nil {
				return fmt.Errorf("create manifest: %w", err)
			}
		}

		return nil
	},
}

func init() {
	initCmd.SilenceUsage = true
	rootCmd.AddCommand(initCmd)
}
