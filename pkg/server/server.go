package server

import (
	"challenge/pkg/config"
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type ChallengeServer struct {
	proto.UnimplementedChallengeServiceServer
}

func Run() {
	logger.Info.Println("Sourcing vars")
	if err := config.Read(".env"); err != nil {
		logger.Error.Fatalln("Error while sourcing vars")
	}
	cfg := config.Get()

	logger.Info.Println("Creating server")
	srv := grpc.NewServer()
	service := ChallengeServer{}
	proto.RegisterChallengeServiceServer(srv, &service)

	connectionString := fmt.Sprintf("localhost:%s", cfg.Port)
	logger.Info.Printf("Starting server on %s", connectionString)
	lis, _ := net.Listen("tcp", connectionString)
	srv.Serve(lis)
}
