package server

import "challenge/pkg/proto"

type ChallengeServer struct {
	proto.UnimplementedChallengeServiceServer
}

func (s *ChallengeServer) StartTimer(timer *proto.Timer, _ proto.ChallengeService_StartTimerServer) error {
	return nil
}
