package gateway

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"os"
	"time"
)

func CreateNewConfig(ctx context.Context) (*SokratesHandler, error) {
	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	var tracer *aristophanes.ClientTracer
	var streamer pb.TraceService_ChorusClient

	maxRetries := 3
	retryDelay := 10 * time.Second

	for i := 1; i <= maxRetries; i++ {
		tracer, err = aristophanes.NewClientTracer(aristophanes.DefaultAddress)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create tracer (attempt %d/%d): %s", i, maxRetries, err.Error()))

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logging.Error("giving up after 3 retries to connect to tracer")
		os.Exit(1)
	}

	for i := 1; i <= maxRetries; i++ {
		streamer, err = tracer.Chorus(ctx)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create chorus streamer (attempt %d/%d): %s", i, maxRetries, err.Error()))
		if i < maxRetries {
			time.Sleep(retryDelay)
		}
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
	multipleChoiceClient, err := triakonta.NewKritiasClient(multipleChoiceClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	multipleChoiceClientHealthy := multipleChoiceClient.WaitForHealthyState()
	if !multipleChoiceClientHealthy {
		logging.Debug("multiplechoice client not ready - restarting seems the only option")
		os.Exit(1)
	}

	authorBasedClientAddress := config.StringFromEnv(config.EnvAuthorBasedClient, config.DefaultAuthorBasedAddress)
	authorBasedClient, err := anabasis.NewXenofonClient(authorBasedClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	authorBasedClientHealthy := authorBasedClient.WaitForHealthyState()
	if !authorBasedClientHealthy {
		logging.Debug("authorbased client not ready - restarting seems the only option")
		os.Exit(1)
	}

	dialogueClientAddress := config.StringFromEnv(config.EnvDialogueClient, config.DefaultDialogueAddress)
	dialogueClient, err := philia.NewKritonClient(dialogueClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	dialogueClientHealthy := dialogueClient.WaitForHealthyState()
	if !dialogueClientHealthy {
		logging.Debug("dialogue client not ready - restarting seems the only option")
		os.Exit(1)
	}

	grammarClientAddress := config.StringFromEnv("ANTISTHENES_SERVICE", config.DefaultGrammarBasedAddress)
	grammarClient, err := kunismos.NewAntisthenesClient(grammarClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	grammarClientHealthy := grammarClient.WaitForHealthyState()
	if !grammarClientHealthy {
		logging.Debug("grammar client not ready - restarting seems the only option")
		os.Exit(1)
	}

	journeyClientAddress := config.StringFromEnv(config.EnvJourneyClient, config.DefaultJourneyAddress)
	journeyClient, err := strategos.NewAlkibiadesClient(journeyClientAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	journeyClientHealthy := journeyClient.WaitForHealthyState()
	if !journeyClientHealthy {
		logging.Debug("grammar client not ready - restarting seems the only option")
		os.Exit(1)
	}

	return &SokratesHandler{
		Streamer:          streamer,
		Randomizer:        randomizer,
		MediaClient:       mediaClient,
		MultiChoiceClient: multipleChoiceClient,
		AuthorBasedClient: authorBasedClient,
		DialogueClient:    dialogueClient,
		GrammarClient:     grammarClient,
		JourneyClient:     journeyClient,
	}, nil
}
