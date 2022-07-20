package client

import (
	"challenge/pkg/config"
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"context"
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var timerName string
var timerDuration int64
var timerFreq int64

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Poll timer",
	Long:  `Poll timer via gRPC`,
	Run: func(_ *cobra.Command, _ []string) {
		// read config to get port
		config.Read(".env")
		cfg := config.Get()

		// listen to server
		logger.Info.Printf("Connecting to localhost:%s\n", cfg.Port)
		conn, err := grpc.Dial("localhost:"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error.Fatalln("Can't create connection for timer client")
		}
		defer conn.Close()

		// create client and stream
		client := proto.NewChallengeServiceClient(conn)
		stream, err := client.StartTimer(context.Background(), &proto.Timer{
			Seconds:   timerDuration,
			Frequency: timerFreq,
			Name:      timerName,
		})
		if err != nil {
			logger.Error.Fatalln("Cannot create remote timer")
		}
		done := make(chan bool)

		go func() {
			for {
				timer, err := stream.Recv()

				if err == io.EOF {
					done <- true
					return
				}

				if err != nil {
					logger.Error.Println(err)
				}

				logger.Info.Printf("%s: %d\n", timer.Name, timer.Seconds)
			}
		}()

		<-done
	},
}

func init() {
	rootCmd.AddCommand(timerCmd)
	timerCmd.Flags().StringVarP(&timerName, "name", "n", "default", "Timer name")
	timerCmd.Flags().Int64VarP(&timerDuration, "seconds", "d", 60, "Timer durations in seconds")
	timerCmd.Flags().Int64VarP(&timerFreq, "freq", "f", 1, "Timer frequency in seconds")
}
