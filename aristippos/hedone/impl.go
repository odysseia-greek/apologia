package hedone

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/aristippos/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type MediaService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
	Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type MediaServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   pbar.TraceService_ChorusClient
	Archytas   archytas.Client
	Progress   *progress.ProgressTracker
	pb.UnimplementedAristipposServer
}

type MediaServiceClient struct {
	Impl MediaService
}

type MediaClient struct {
	media pb.AristipposClient
}

func NewAristipposClient(address string) (*MediaClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewAristipposClient(conn)
	return &MediaClient{media: client}, nil
}

func (m *MediaClient) WaitForHealthyState() bool {
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

func (m *MediaClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return m.media.Health(ctx, request)
}

func (m *MediaClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return m.media.Options(ctx, request)
}

func (m *MediaClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return m.media.Question(ctx, request)
}

func (m *MediaClient) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error) {
	return m.media.Answer(ctx, request)
}
