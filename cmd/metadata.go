package client

import (
	"challenge/pkg/logger"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "meta",
	Short: "Get metadata",
	Long:  `Get metadata via gRPC`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.InfoLogger.Println("Reading metadata...")
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
