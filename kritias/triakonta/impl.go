package triakonta

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/kritias/gen/go/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MultipleChoiceService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *v1.OptionsRequest) (*v1.AggregatedOptions, error)
	Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error)
	Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type MultipleChoiceServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   arv1.TraceService_ChorusClient
	Archytas   archytas.Client
	Progress   *progress.ProgressTracker
	v1.UnimplementedKritiasServer
}

type MultipleChoiceServiceClient struct {
	Impl MultipleChoiceService
}

type MutpleChoiceClient struct {
	multiplechoice v1.KritiasClient
}

func NewKritiasClient(address string) (*MutpleChoiceClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewKritiasClient(conn)
	return &MutpleChoiceClient{multiplechoice: client}, nil
}

func (m *MutpleChoiceClient) WaitForHealthyState() bool {
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

func (m *MutpleChoiceClient) Health(ctx context.Context, request *v1.HealthRequest) (*v1.HealthResponse, error) {
	return m.multiplechoice.Health(ctx, request)
}

func (m *MutpleChoiceClient) Options(ctx context.Context, request *v1.OptionsRequest) (*v1.AggregatedOptions, error) {
	return m.multiplechoice.Options(ctx, request)
}

func (m *MutpleChoiceClient) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	return m.multiplechoice.Question(ctx, request)
}

func (m *MutpleChoiceClient) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	return m.multiplechoice.Answer(ctx, request)
}
