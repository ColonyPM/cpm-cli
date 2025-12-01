package cmd

import (
	"os"

	"github.com/ColonyPM/cpm-cli/internal/repository"
	"github.com/spf13/cobra"
)

const baseURL = "https://conjoined-aide-abeyantly.ngrok-free.dev/api"

// Shared client for all commands in this package.
var httpClient = repository.New(baseURL)

// var defaultDestPath = "~/.local/share/cpm/pkgs"
var defaultDestPath = "./pkgs"

var rootCmd = &cobra.Command{
	Use:   "cpm-cli",
	Short: "A brief description of your application",
	Long:  "A longer description that spans multiple lines",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
