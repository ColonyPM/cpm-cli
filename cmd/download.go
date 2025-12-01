package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a package from the repository",
	Long:  "Download a package from the repository - download <package-name>",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("download called")

		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		if err := httpClient.Download(ctx, args[0], defaultDestPath); err != nil {
			return err
		}

		fmt.Println("Download OK")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
