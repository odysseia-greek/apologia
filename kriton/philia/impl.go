package philia

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/kriton/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type DialogueService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
	Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error)
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
	Streamer   pbar.TraceService_ChorusClient
	Archytas   archytas.Client
	pb.UnimplementedKritonServer
}

type DialogueServiceClient struct {
	Impl DialogueService
}
type DialogueClient struct {
	dialogue pb.KritonClient
}

func NewKritonClient(address string) (*DialogueClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewKritonClient(conn)
	return &DialogueClient{dialogue: client}, nil
}

func (d *DialogueClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := d.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (d *DialogueClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return d.dialogue.Health(ctx, request)
}

func (d *DialogueClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return d.dialogue.Options(ctx, request)
}

func (d *DialogueClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return d.dialogue.Question(ctx, request)
}

func (d *DialogueClient) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	return d.dialogue.Answer(ctx, request)
}
