package server

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestReadMetadata(t *testing.T) {
	s := ChallengeServer{}
	tests := []struct {
		value       string
		shouldThrow bool
	}{
		{
			value:       "foo",
			shouldThrow: false,
		},
		{
			value:       "bar",
			shouldThrow: false,
		},
		{
			value:       "",
			shouldThrow: true,
		},
	}
	for _, tc := range tests {
		resp, err := s.ReadMetadata(metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
			DesiredKey: tc.value})), nil)
		if err != nil && !tc.shouldThrow {
			t.Error("Unexpected exception while reading metadata")
		}
		if resp.Data != tc.value {
			t.Errorf("Wanted %s. Got %s", tc.value, resp.GetData())
		}
	}
}
