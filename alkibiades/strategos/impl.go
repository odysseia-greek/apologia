package strategos

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/alkibiades/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type JourneyService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
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
	Streamer   pbar.TraceService_ChorusClient
	Archytas   archytas.Client
	pb.UnimplementedAlkibiadesServer
}

type JourneyServiceClient struct {
	Impl JourneyService
}

type JourneyClient struct {
	journey pb.AlkibiadesClient
}

func NewAlkibiadesClient(address string) (*JourneyClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewAlkibiadesClient(conn)
	return &JourneyClient{journey: client}, nil
}

func (j *JourneyClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := j.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (j *JourneyClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return j.journey.Health(ctx, request)
}

func (j *JourneyClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return j.journey.Options(ctx, request)
}
func (j *JourneyClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return j.journey.Question(ctx, request)
}
