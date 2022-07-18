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

func (s *ChallengeServer) StartTimer(timer *proto.Timer, _ proto.ChallengeService_StartTimerServer) error {
	return nil
}

func Run() {
	logger.InfoLogger.Println("Sourcing vars")
	config.Read(".env")
	cfg := config.Get()

	logger.InfoLogger.Println("Creating server")
	srv := grpc.NewServer()
	service := ChallengeServer{}
	proto.RegisterChallengeServiceServer(srv, &service)

	connectionString := fmt.Sprintf("localhost:%s", cfg.Port)
	logger.InfoLogger.Printf("Starting server on %s", connectionString)
	lis, _ := net.Listen("tcp", connectionString)
	srv.Serve(lis)
}
