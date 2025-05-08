package kunismos

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/progress"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/apologia/antisthenes/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type GrammarService interface {
	WaitForHealthyState() bool
	Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error)
	Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error)
	Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type GrammarServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   pbar.TraceService_ChorusClient
	Archytas   archytas.Client
	Progress   *progress.ProgressTracker
	pb.UnimplementedAntisthenesServer
}

type GrammarServiceClient struct {
	Impl GrammarService
}

type GrammarClient struct {
	grammar pb.AntisthenesClient
}

func NewAntisthenesClient(address string) (*GrammarClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewAntisthenesClient(conn)
	return &GrammarClient{grammar: client}, nil
}

func (g *GrammarClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := g.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (g *GrammarClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return g.grammar.Health(ctx, request)
}

func (g *GrammarClient) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	return g.grammar.Options(ctx, request)
}

func (g *GrammarClient) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	return g.grammar.Question(ctx, request)
}

func (g *GrammarClient) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error) {
	return g.grammar.Answer(ctx, request)
}
