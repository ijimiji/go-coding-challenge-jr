package client

import (
	"challenge/pkg/config"
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var mdValue string

var metadataCmd = &cobra.Command{
	Use:   "meta",
	Short: "Get metadata",
	Long:  `Get metadata via gRPC`,
	Run: func(_ *cobra.Command, _ []string) {
		logger.Info.Println("Reading metadata...")

		config.Read(".env")
		cfg := config.Get()

		logger.Info.Printf("Connecting to localhost:%s\n", cfg.Port)
		conn, err := grpc.Dial("localhost:"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			logger.Error.Fatalln("Can't create connection for metadata client")
		}

		defer conn.Close()
		client := proto.NewChallengeServiceClient(conn)
		value, err := client.ReadMetadata(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"i-am-random-key": mdValue})), &proto.Placeholder{Data: "foo"})

		if err != nil {
			logger.Error.Fatalln(err.Error())
		}

		fmt.Printf("Extracted metadata: %s\n", value.GetData())

	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
	metadataCmd.Flags().StringVarP(&mdValue, "metadata", "m", "", "Metadata to be extracted")
}
