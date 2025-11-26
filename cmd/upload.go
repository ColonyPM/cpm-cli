package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a package to the repository",
	Long:  "cpm upload <dir> <token> \n\n  A token is generated from the ColonyPM web server.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Uploading package from directory %s to repository with token %s...\n", args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
