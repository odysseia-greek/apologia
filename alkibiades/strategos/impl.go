package strategos

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/alkibiades/gen/go/v1"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JourneyService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error)
	Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type JourneyServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   arv1.TraceService_ChorusClient
	Archytas   archytas.Client
	v1.UnimplementedAlkibiadesServer
}

type JourneyServiceClient struct {
	Impl JourneyService
}

type JourneyClient struct {
	journey v1.AlkibiadesClient
}

func NewAlkibiadesClient(address string) (*JourneyClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewAlkibiadesClient(conn)
	return &JourneyClient{journey: client}, nil
}

func (j *JourneyClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := j.Health(context.Background(), &koinosv1.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (j *JourneyClient) Health(ctx context.Context, request *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	return j.journey.Health(ctx, request)
}

func (j *JourneyClient) Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error) {
	return j.journey.Options(ctx, request)
}

func (j *JourneyClient) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	return j.journey.Question(ctx, request)
}
