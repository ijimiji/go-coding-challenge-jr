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

// re-use this key name in other parts of program
var DesiredKey = "i-am-random-key"

// take metadata from incoming context and return it in placeholder
// blank metadata is unacceptable
func (s *ChallengeServer) ReadMetadata(ctx context.Context, placeHolder *proto.Placeholder) (*proto.Placeholder, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error.Println("Can't read metadata from context")
		return nil, errors.New("can't read metadata from context")
	}

	data := ""
	if arrData := md.Get(DesiredKey); len(arrData) > 0 {
		data = arrData[0]
	}

	if data == "" {
		logger.Error.Println("Can't read metadata from context")
		return nil, status.Errorf(codes.InvalidArgument, "no key/value was provided")
	}

	return &proto.Placeholder{Data: data}, nil
}
