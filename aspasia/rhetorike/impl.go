package rhetorike

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/hesiodos"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GathererService interface {
	WaitForHealthyState() bool
	Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type GathererServiceImpl struct {
	Version     string
	Archytas    archytas.Client
	Client      service.OdysseiaClient
	FuzzyClient *hesiodos.GenericGrpcClient[*monophthalmus.FuzzyClient]
	Streamer    arv1.TraceService_ChorusClient
	v1.UnimplementedAspasiaServiceServer
}

type GathererServiceClient struct {
	Impl GathererService
}
type GathererClient struct {
	gatherer v1.AspasiaServiceClient
}

func NewAspasiaClient(address string) (*GathererClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewAspasiaServiceClient(conn)
	return &GathererClient{gatherer: client}, nil
}

func (g *GathererClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		return true
	}

	return false
}

func (g *GathererClient) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	return g.gatherer.Search(ctx, request)
}
