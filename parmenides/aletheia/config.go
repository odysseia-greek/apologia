package aletheia

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	eupalinos "github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	aristarchos "github.com/odysseia-greek/alexandreia/aristarchos/scholar"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pbp "github.com/odysseia-greek/delphi/aristides/proto"
	"google.golang.org/grpc/metadata"
)

func CreateNewConfig() (*ParmenidesHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)
	if err != nil {
		return nil, err
	}

	healthy := ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceId := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	md := metadata.New(map[string]string{service.HeaderKey: traceId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ctx, &pbp.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	elasticService := aristoteles.ElasticService(tls)

	cfg = models.Config{
		Service:     elasticService,
		Username:    vaultConfig.ElasticUsername,
		Password:    vaultConfig.ElasticPassword,
		ElasticCERT: vaultConfig.ElasticCERT,
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	channel := config.StringFromEnv(config.EnvChannel, config.DefaultParmenidesChannel)

	index := config.StringFromEnv(config.EnvIndex, "")
	if index == "" {
		return nil, fmt.Errorf("no index found in environment please set %s", config.EnvIndex)
	}

	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)
	logging.Debug(fmt.Sprintf("creating new eupalinos client: %s", eupalinosAddress))
	queue, err := eupalinos.NewEupalinosClient(eupalinosAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	logging.Debug(fmt.Sprintf("created new eupalinos client with channel and dutch: %s + %s + %s", eupalinosAddress, channel, config.DefaultDutchChannel))
	logging.Debug("waiting for queue to be ready")
	queueHealthy := queue.WaitForHealthyState()
	if !queueHealthy {
		logging.Debug("no queue that is healthy")
	}

	policyName := config.StringFromEnv("HOT_POLICY_NAME", "hot_plain")

	handler := &ParmenidesHandler{
		Index:            index,
		Created:          0,
		Elastic:          elastic,
		Eupalinos:        queue,
		Channel:          channel,
		DutchChannel:     config.DefaultDutchChannel,
		PolicyName:       policyName,
		Ambassador:       ambassador,
		Aggregator:       nil,
		AggregatorCancel: nil,
	}

	if index == "author-based-quiz" || index == "grammar-quiz" {
		aggregatorAddress := config.StringFromEnv(config.EnvAggregatorAddress, config.DefaultAggregatorAddress)
		aggregator, err := aristarchos.NewClientAggregator(aggregatorAddress)
		if err != nil {
			logging.Error(err.Error())
			return nil, err
		}

		logging.Debug(fmt.Sprintf("creating new aggregator client: %s", aggregatorAddress))
		logging.Debug("waiting for aggregator to be ready")
		aggregatorHealthy := aggregator.WaitForHealthyState()
		if !aggregatorHealthy {
			logging.Debug("aggregator service not ready - restarting seems the only option")
			os.Exit(1)
		}

		logging.Debug("aggregator is ready")
		// New context for aggregator streamer
		aggrContext, aggregatorCancel := context.WithCancel(context.Background())
		aristarchosStreamer, err := aggregator.CreateNewEntry(aggrContext)
		if err != nil {
			logging.Error(err.Error())
			return nil, err
		}

		handler.Aggregator = aristarchosStreamer
		handler.AggregatorCancel = aggregatorCancel
	}

	return handler, nil
}
