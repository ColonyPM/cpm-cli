package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a package",
	Long:  "install <package-name> - Install a package from the repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Installing package %s...", args[0])
		fmt.Printf("Deploying package '%s'...\n on colony...", args[0])
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
