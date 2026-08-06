package rhetorike

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
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
	Version           string
	AlexandrosAddress string
	GraphqlClient     *http.Client
	Archytas          archytas.Client
	Dionysios         DionysiosResearchClient
	DionysiosConn     *grpc.ClientConn
	Streamer          arv1.TraceService_ChorusClient
	v1.UnimplementedAspasiaServiceServer
}

type DionysiosResearchClient interface {
	Research(ctx context.Context, request *dionysiosv1.ResearchRequest, opts ...grpc.CallOption) (*dionysiosv1.ResearchResponse, error)
}

func (g *GathererServiceImpl) Close() error {
	if g.DionysiosConn != nil {
		if err := g.DionysiosConn.Close(); err != nil {
			return err
		}
	}
	if g.Archytas != nil {
		return g.Archytas.Close()
	}
	return nil
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
