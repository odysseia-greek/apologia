package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/aspasia/rhetorike"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

func CreateNewConfig(ctx context.Context) (*SokratesHandler, error) {
	start := time.Now()

	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	var tracer *aristophanes.ClientTracer
	var streamer arv1.TraceService_ChorusClient

	maxRetries := 10
	retryDelay := 3 * time.Second

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

	healthyTracer := false
	if tracer != nil {
		healthyTracer = tracer.WaitForHealthyState()
	}

	mediaClientAddress := config.StringFromEnv(config.EnvMediaClient, config.DefaultMediaAddress)
	mediaClient, err := NewGenericGrpcClient[*hedone.MediaClient](
		mediaClientAddress,
		hedone.NewAristipposClient,
	)
	if err != nil {
		return nil, err
	}

	mediaClientHealthy := false
	if mediaClient != nil {
		mediaClientHealthy = mediaClient.client.WaitForHealthyState()
	}

	multipleChoiceClientAddress := config.StringFromEnv(config.EnvMultiChoiceClient, config.DefaultMultiChoiceAddress)
	multipleChoiceClient, err := NewGenericGrpcClient[*triakonta.MutpleChoiceClient](
		multipleChoiceClientAddress,
		triakonta.NewKritiasClient,
	)
	if err != nil {
		return nil, err
	}

	multipleChoiceClientHealthy := false
	if multipleChoiceClient != nil {
		multipleChoiceClientHealthy = multipleChoiceClient.client.WaitForHealthyState()
	}

	authorBasedClientAddress := config.StringFromEnv(config.EnvAuthorBasedClient, config.DefaultAuthorBasedAddress)
	authorBasedClient, err := NewGenericGrpcClient[*anabasis.AuthorBasedClient](
		authorBasedClientAddress,
		anabasis.NewXenofonClient,
	)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	authorBasedClientHealthy := false
	if authorBasedClient != nil {
		authorBasedClientHealthy = authorBasedClient.client.WaitForHealthyState()
	}

	dialogueClientAddress := config.StringFromEnv(config.EnvDialogueClient, config.DefaultDialogueAddress)
	dialogueClient, err := NewGenericGrpcClient[*philia.DialogueClient](
		dialogueClientAddress,
		philia.NewKritonClient,
	)

	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	dialogueClientHealthy := false
	if dialogueClient != nil {
		dialogueClientHealthy = dialogueClient.client.WaitForHealthyState()
	}

	grammarClientAddress := config.StringFromEnv(config.EnvGrammarBasedClient, config.DefaultGrammarBasedAddress)
	grammarClient, err := NewGenericGrpcClient[*kunismos.GrammarClient](
		grammarClientAddress,
		kunismos.NewAntisthenesClient,
	)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	grammarClientHealthy := false
	if grammarClient != nil {
		grammarClientHealthy = grammarClient.client.WaitForHealthyState()
	}

	journeyClientAddress := config.StringFromEnv(config.EnvJourneyClient, config.DefaultJourneyAddress)
	journeyClient, err := NewGenericGrpcClient[*strategos.JourneyClient](
		journeyClientAddress,
		strategos.NewAlkibiadesClient,
	)

	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	journeyClientHealthy := false
	if journeyClient != nil {
		journeyClientHealthy = journeyClient.client.WaitForHealthyState()
	}

	gathererClientAddress := config.StringFromEnv("ASPASIA_SERVICE", "aspasia:50060")
	gathererClient, err := NewGenericGrpcClient[*rhetorike.GathererClient](
		gathererClientAddress,
		rhetorike.NewAspasiaClient,
	)

	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	gathererClientHealthy := false
	if gathererClient != nil {
		gathererClientHealthy = gathererClient.client.WaitForHealthyState()
	}

	elapsed := time.Since(start)

	logging.System(fmt.Sprintf(`Sokrates Configuration Overview:
- Initialization Time: %s
- Tracer Service:      %v (Address: %s)
- Aristipoos Service:   %v (Address: %s)
- Kritias Service:   %v (Address: %s)
- Xenofon Service:  %v (Address: %s)
- Kriton Service:  %v (Address: %s)
- Antisthenes Service:   %v (Address: %s)
- Alkibiades Service:  %v (Address: %s)
- Aspasia Service:   %v (Address: %s)
`,
		elapsed,
		healthyTracer, aristophanes.DefaultAddress,
		mediaClientHealthy, mediaClientAddress,
		multipleChoiceClientHealthy, multipleChoiceClientAddress,
		authorBasedClientHealthy, authorBasedClientAddress,
		dialogueClientHealthy, dialogueClientAddress,
		grammarClientHealthy, grammarClientAddress,
		journeyClientHealthy, journeyClientAddress,
		gathererClientHealthy, gathererClientAddress,
	))

	return &SokratesHandler{
		Streamer:          streamer,
		Randomizer:        randomizer,
		MediaClient:       mediaClient,
		MultiChoiceClient: multipleChoiceClient,
		AuthorBasedClient: authorBasedClient,
		DialogueClient:    dialogueClient,
		GrammarClient:     grammarClient,
		JourneyClient:     journeyClient,
		GathererClient:    gathererClient,
	}, nil
}
