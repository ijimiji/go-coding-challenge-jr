// This part of server module is responsible for extracting metadata from contexts
package server

import (
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var desiredKey = "i-am-random-key"

func (s *ChallengeServer) ReadMetadata(ctx context.Context, placeHolder *proto.Placeholder) (*proto.Placeholder, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.ErrorLogger.Println("Can't read metadata from context")
		return nil, errors.New("can't read metadata from context")
	}

	data := ""
	if arrData := md.Get(desiredKey); len(arrData) > 0 {
		data = arrData[0]
	}

	if data == "" {
		logger.ErrorLogger.Println("Can't read metadata from context")
		return nil, status.Errorf(codes.InvalidArgument, "no i-am-random-key was provided")
	}

	return &proto.Placeholder{Data: data}, nil
}
