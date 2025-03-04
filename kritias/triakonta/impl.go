package triakonta

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/kritias/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type MultipleChoiceService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
	Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type MultipleChoiceServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   pbar.TraceService_ChorusClient
	pb.UnimplementedKritiasServer
}

type MultipleChoiceServiceClient struct {
	Impl MultipleChoiceService
}

type MutpleChoiceClient struct {
	multiplechoice pb.KritiasClient
}

func NewAristipposClient(address string) (*MutpleChoiceClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewKritiasClient(conn)
	return &MutpleChoiceClient{multiplechoice: client}, nil
}

func (m *MutpleChoiceClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MutpleChoiceClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return m.multiplechoice.Health(ctx, request)
}

func (m *MutpleChoiceClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return m.multiplechoice.Options(ctx, request)
}

func (m *MutpleChoiceClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return m.multiplechoice.Question(ctx, request)
}

func (m *MutpleChoiceClient) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error) {
	return m.multiplechoice.Answer(ctx, request)
}
