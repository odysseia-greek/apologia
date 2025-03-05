package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"os"
)

func CreateNewConfig(ctx context.Context) (*SokratesHandler, error) {
	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - starting up without traces")
	}

	mediaClientAddress := config.StringFromEnv(config.EnvMediaClient, config.DefaultMediaAddress)
	mediaClient, err := hedone.NewAristipposClient(mediaClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	mediaClientHealthy := mediaClient.WaitForHealthyState()
	if !mediaClientHealthy {
		logging.Debug("media client not ready - restarting seems the only option")
		os.Exit(1)
	}

	multipleChoiceClientAddress := config.StringFromEnv(config.EnvMultiChoiceClient, config.DefaultMultiChoiceAddress)
	multipleChoiceClient, err := triakonta.NewAristipposClient(multipleChoiceClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	multipleChoiceClientHealthy := multipleChoiceClient.WaitForHealthyState()
	if !multipleChoiceClientHealthy {
		logging.Debug("media client not ready - restarting seems the only option")
		os.Exit(1)
	}

	return &SokratesHandler{
		Streamer:          streamer,
		Randomizer:        randomizer,
		MediaClient:       mediaClient,
		MultiChoiceClient: multipleChoiceClient,
	}, nil
}
