package anabasis

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/xenofon/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthorBasedService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
	Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error)
	WordForms(ctx context.Context, request *pb.WordFormRequest) (*pb.WordFormRequest, error)
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
	Streamer   pbar.TraceService_ChorusClient
	pb.UnimplementedXenofonServer
}
type AuthorBasedServiceClient struct {
	Impl AuthorBasedService
}
type AuthorBasedClient struct {
	authorbased pb.XenofonClient
}

func NewXenofonClient(address string) (*AuthorBasedClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewXenofonClient(conn)
	return &AuthorBasedClient{authorbased: client}, nil
}

func (m *AuthorBasedClient) WaitForHealthyState() bool {
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

func (m *AuthorBasedClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return m.authorbased.Health(ctx, request)
}

func (m *AuthorBasedClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return m.authorbased.Options(ctx, request)
}

func (m *AuthorBasedClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return m.authorbased.Question(ctx, request)
}

func (m *AuthorBasedClient) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	return m.authorbased.Answer(ctx, request)
}

func (m *AuthorBasedClient) WordForms(ctx context.Context, request *pb.WordFormRequest) (*pb.WordFormResponse, error) {
	return m.authorbased.WordForms(ctx, request)
}
