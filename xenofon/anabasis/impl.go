package anabasis

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/apologia/xenofon/gen/go/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthorBasedService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *koinosv1.OptionsRequest) (*koinosv1.AggregatedOptions, error)
	Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error)
	Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error)
	WordForms(ctx context.Context, request *v1.WordFormRequest) (*v1.WordFormRequest, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type AuthorBasedServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Archytas   archytas.Client
	Progress   *progress.ProgressTracker
	Streamer   arv1.TraceService_ChorusClient
	v1.UnimplementedXenofonServer
}
type AuthorBasedServiceClient struct {
	Impl AuthorBasedService
}
type AuthorBasedClient struct {
	authorbased v1.XenofonClient
}

func NewXenofonClient(address string) (*AuthorBasedClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewXenofonClient(conn)
	return &AuthorBasedClient{authorbased: client}, nil
}

func (m *AuthorBasedClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.Health(context.Background(), &koinosv1.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *AuthorBasedClient) Health(ctx context.Context, request *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	return m.authorbased.Health(ctx, request)
}

func (m *AuthorBasedClient) Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error) {
	return m.authorbased.Options(ctx, request)
}

func (m *AuthorBasedClient) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	return m.authorbased.Question(ctx, request)
}

func (m *AuthorBasedClient) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	return m.authorbased.Answer(ctx, request)
}

func (m *AuthorBasedClient) WordForms(ctx context.Context, request *v1.WordFormRequest) (*v1.WordFormResponse, error) {
	return m.authorbased.WordForms(ctx, request)
}
