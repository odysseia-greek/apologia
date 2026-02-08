package philia

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/apologia/kriton/gen/go/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DialogueService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error)
	Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error)
	Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type DialogueServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   arv1.TraceService_ChorusClient
	Archytas   archytas.Client
	v1.UnimplementedKritonServer
}

type DialogueServiceClient struct {
	Impl DialogueService
}
type DialogueClient struct {
	dialogue v1.KritonClient
}

func NewKritonClient(address string) (*DialogueClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewKritonClient(conn)
	return &DialogueClient{dialogue: client}, nil
}

func (d *DialogueClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := d.Health(context.Background(), &koinosv1.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (d *DialogueClient) Health(ctx context.Context, request *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	return d.dialogue.Health(ctx, request)
}

func (d *DialogueClient) Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error) {
	return d.dialogue.Options(ctx, request)
}

func (d *DialogueClient) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	return d.dialogue.Question(ctx, request)
}

func (d *DialogueClient) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	return d.dialogue.Answer(ctx, request)
}
