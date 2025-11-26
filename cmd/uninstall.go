package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "cpm uninstall <package-name> - Uninstall a package",
	Long:  "cpm uninstall <package-name> - Uninstall a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Uninstalling package %s...\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
