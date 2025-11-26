package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all installed packages",
	Long:  "list - Lists all installed packages in the ColonyOS environment",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all installed packages: \n- packageA\n- packageB\n- packageC")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
