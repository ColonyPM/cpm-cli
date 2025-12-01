package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const manifestName = "package.yml"

const manifestTemplate = `name: "%s"
version: "0.0.0"
description: ""
author: "%s"
`

var ErrManifestAlreadyExists = errors.New("manifest already exists")

func getUsername() (string, error) {
	// Preferred: os/user
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username, nil
	}

	// Fallbacks via environment
	if v := os.Getenv("USER"); v != "" {
		return v, nil
	}
	if v := os.Getenv("USERNAME"); v != "" {
		return v, nil
	}

	return "", errors.New("failed to get username")
}

func initPackage() error {
	// Fail if package.yml already exists
	if _, err := os.Stat(manifestName); err == nil {
		return ErrManifestAlreadyExists
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking %s: %v", manifestName, err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get cwd: %v", err)
	}

	// Use the cwd name as the package name
	pkgName := filepath.Base(cwd)

	username, err := getUsername()
	if err != nil {
		username = ""
	}

	escapedPkg := strings.ReplaceAll(pkgName, `"`, `\"`)
	escapedUser := strings.ReplaceAll(username, `"`, `\"`)

	content := fmt.Sprintf(manifestTemplate, escapedPkg, escapedUser)

	// perms: rw-r--r--
	if err := os.WriteFile(manifestName, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %v", manifestName, err)
	}

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init [DIR]",
	Short: "Initialize a new package",
	Long:  "Initialize a new cpm package in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initPackage()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
