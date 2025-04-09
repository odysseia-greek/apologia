package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

type SokratesHandler struct {
	Streamer          pbar.TraceService_ChorusClient
	Randomizer        randomizer.Random
	GrammarClient     *kunismos.GrammarClient
	MediaClient       *hedone.MediaClient
	MultiChoiceClient *triakonta.MutpleChoiceClient
	AuthorBasedClient *anabasis.AuthorBasedClient
	DialogueClient    *philia.DialogueClient
}

func (s *SokratesHandler) createRequestHeader(requestID, sessionId string) (context.Context, context.CancelFunc) {
	requestCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	md := metadata.New(map[string]string{config.HeaderKey: requestID,
		config.SessionIdKey: sessionId})
	requestCtx = metadata.NewOutgoingContext(requestCtx, md)

	return requestCtx, ctxCancel
}

func (s *SokratesHandler) Health(requestID, sessionId string) (*model.AggregatedHealthResponse, error) {
	var services []*model.ServiceHealth
	allHealthy := true

	checkService := func(name string, client func(ctx context.Context, in *pbartrippos.HealthRequest) (*pbartrippos.HealthResponse, error)) {
		healthCtx, cancel := s.createRequestHeader(requestID, sessionId)
		defer cancel()

		resp, err := client(healthCtx, &pbartrippos.HealthRequest{})
		serviceHealth := &model.ServiceHealth{
			Name:           &name,
			Healthy:        &resp.Healthy,
			DatabaseHealth: &resp.DatabaseHealth.Healthy,
			Version:        &resp.Version,
		}

		if err != nil {
			allHealthy = false
		} else {
			serviceHealth.Healthy = &resp.Healthy
			serviceHealth.DatabaseHealth = &resp.DatabaseHealth.Healthy
			serviceHealth.Version = &resp.Version
			if !resp.Healthy {
				allHealthy = false
			}
		}

		services = append(services, serviceHealth)
	}

	checkService("media", s.MediaClient.Health)

	return &model.AggregatedHealthResponse{
		Healthy:  &allHealthy,
		Time:     ptr(time.Now().Format(time.RFC3339)),
		Version:  ptr(os.Getenv("VERSION")),
		Services: services,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}
