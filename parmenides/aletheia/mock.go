package aletheia

import (
	"context"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
)

// MockEupalinosClient is a mock implementation of EupalinosClient
type MockEupalinosClient struct {
}

// EnqueueMessage is the mock implementation for the EnqueueMessage method
func (m *MockEupalinosClient) EnqueueMessage(ctx context.Context, in *pb.Epistello, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	// Return a mock response or an error based on your test scenario
	// For example:
	return &pb.EnqueueResponse{
		Id: "your-generated-uuid",
	}, nil
}
