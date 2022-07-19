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
)

var longUrl string
var shortenCmd = &cobra.Command{
	Use:   "link",
	Short: "Shorten link",
	Long:  `Shorten link via gRPC`,
	Run: func(_ *cobra.Command, _ []string) {

		config.Read(".env")
		cfg := config.Get()

		if longUrl == "" {
			logger.Error.Fatalln("No link was provided")
		}

		logger.Info.Printf("Connecting to localhost:%s\n", cfg.Port)
		conn, err := grpc.Dial("localhost:"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error.Fatalln("Can't create connection for shortener client")
		}
		defer conn.Close()
		client := proto.NewChallengeServiceClient(conn)
		link, err := client.MakeShortLink(context.Background(), &proto.Link{Data: longUrl})
		if err != nil {
			logger.Error.Fatalln(err.Error())
		}

		fmt.Printf("Shortened link: %s\n", link.GetData())
	},
}

func init() {
	rootCmd.AddCommand(shortenCmd)
	shortenCmd.Flags().StringVarP(&longUrl, "link", "l", "", "Link to be shortened")
}
