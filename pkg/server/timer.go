package server

import "challenge/pkg/proto"

func (s *ChallengeServer) StartTimer(timer *proto.Timer, _ proto.ChallengeService_StartTimerServer) error {
	return nil
}
