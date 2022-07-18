package client

import (
	"challenge/pkg/logger"

	"github.com/spf13/cobra"
)

var shortenCmd = &cobra.Command{
	Use:   "link",
	Short: "Shorten link",
	Long:  `Shorten link via gRPC`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.InfoLogger.Println("Shortening url...")
	},
}

func init() {
	rootCmd.AddCommand(shortenCmd)
}
