package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cpm",
	Short: "cpm is a package manager for ColonyOS",
	Long:  "cpm is a package manager for ColonyOS - install, update, and manage packages with ease.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintf(os.Stderr, "Oops. An error while executing cpm '%s'\n", err)
		os.Exit(1)
	}
}
