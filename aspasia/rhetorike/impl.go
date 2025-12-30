package rhetorike

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	koinos "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ExtendedService interface {
	WaitForHealthyState() bool
	ExtendedSearch(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type ExtendedServiceImpl struct {
	Version  string
	Archytas archytas.Client
	v1.UnimplementedAspasiaServiceServer
}

type ExtendedServiceClient struct {
	Impl ExtendedService
}
type ExtendedClient struct {
	extended v1.PtolemaiosServiceClient
}

func NewAspasiaClient(address string) (*ExtendedClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewAspasiaServiceClient(conn)
	return &ExtendedClient{extended: client}, nil
}

func (e *ExtendedClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := e.Health(context.Background(), &emptypb.Empty{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (e *ExtendedClient) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	return e.extended.Health(ctx, request)
}

func (e *ExtendedClient) ExtendedSearch(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	return e.extended.Search(ctx, request)
}
