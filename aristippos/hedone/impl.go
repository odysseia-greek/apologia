package hedone

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aristippos/gen/go/v1"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type MediaService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *v1.OptionsRequest) (*v1.AggregatedOptions, error)
	Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error)
	Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error)
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
	v1.UnimplementedAristipposServer
}

type MediaServiceClient struct {
	Impl MediaService
}

type MediaClient struct {
	media v1.AristipposClient
}

func NewAristipposClient(address string) (*MediaClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewAristipposClient(conn)
	return &MediaClient{media: client}, nil
}

func (m *MediaClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.Health(context.Background(), &v1.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MediaClient) Health(ctx context.Context, request *v1.HealthRequest) (*v1.HealthResponse, error) {
	return m.media.Health(ctx, request)
}

func (m *MediaClient) Options(ctx context.Context, request *v1.OptionsRequest) (*v1.AggregatedOptions, error) {
	return m.media.Options(ctx, request)
}

func (m *MediaClient) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	return m.media.Question(ctx, request)
}

func (m *MediaClient) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	return m.media.Answer(ctx, request)
}
