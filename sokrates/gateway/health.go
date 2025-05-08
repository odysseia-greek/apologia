package gateway

import (
	"context"
	"os"
	"time"

	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"

	pbalkibiades "github.com/odysseia-greek/apologia/alkibiades/proto"
	pbantisthenes "github.com/odysseia-greek/apologia/antisthenes/proto"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	pbkriton "github.com/odysseia-greek/apologia/kriton/proto"
	pbxenofon "github.com/odysseia-greek/apologia/xenofon/proto"

	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) Health(requestID, sessionId string) (*model.AggregatedHealthResponse, error) {
	var services []*model.ServiceHealth
	allHealthy := true

	type healthCheck struct {
		name   string
		client func(ctx context.Context) (bool, *model.DatabaseInfo, *string) // healthy, dbHealthy, version
	}

	checks := []healthCheck{
		{
			name: "media",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbartrippos.HealthResponse
				err := s.MediaClient.CallWithReconnect(func(c *hedone.MediaClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbartrippos.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "multiple-choice",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbkritias.HealthResponse
				err := s.MultiChoiceClient.CallWithReconnect(func(c *triakonta.MutpleChoiceClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbkritias.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "author-based",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbxenofon.HealthResponse
				err := s.AuthorBasedClient.CallWithReconnect(func(c *anabasis.AuthorBasedClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbxenofon.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "dialogue",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbkriton.HealthResponse
				err := s.DialogueClient.CallWithReconnect(func(c *philia.DialogueClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbkriton.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "grammar",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbantisthenes.HealthResponse
				err := s.GrammarClient.CallWithReconnect(func(c *kunismos.GrammarClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbantisthenes.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "journey",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *pbalkibiades.HealthResponse
				err := s.JourneyClient.CallWithReconnect(func(c *strategos.JourneyClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &pbalkibiades.HealthRequest{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}

				databaseHealth := &model.DatabaseInfo{
					Healthy:       &resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
	}

	for _, check := range checks {
		healthCtx, cancel := s.createRequestHeader(requestID, sessionId)
		healthy, dbHealthy, version := check.client(healthCtx)
		cancel()

		serviceHealth := &model.ServiceHealth{
			Name:         &check.name,
			Healthy:      &healthy,
			DatabaseInfo: dbHealthy,
			Version:      version,
		}

		if !healthy {
			allHealthy = false
		}

		services = append(services, serviceHealth)
	}

	return &model.AggregatedHealthResponse{
		Healthy:  ptr(allHealthy),
		Time:     ptr(time.Now().Format(time.RFC3339)),
		Version:  ptr(os.Getenv("VERSION")),
		Services: services,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}
