package aletheia

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pbp "github.com/odysseia-greek/delphi/aristides/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

type EupalinosClient interface {
	EnqueueMessage(ctx context.Context, in *pb.Epistello, opts ...grpc.CallOption) (*pb.EnqueueResponse, error)
}

func CreateNewConfig() (*ParmenidesHandler, *grpc.ClientConn, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
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
		return nil, nil, err
	}

	channel := config.StringFromEnv(config.EnvChannel, config.DefaultParmenidesChannel)
	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)

	index := config.StringFromEnv(config.EnvIndex, "")
	if index == "" {
		return nil, nil, fmt.Errorf("no index found in environment please set %s", config.EnvIndex)
	}

	client, conn, err := createEupalinosClient(eupalinosAddress)
	if err != nil {
		return nil, nil, err
	}

	policyName := fmt.Sprintf("%s_policy", index)

	handler := &ParmenidesHandler{
		Index:            index,
		Created:          0,
		Elastic:          elastic,
		Eupalinos:        client,
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
			return nil, nil, err
		}
		aggregatorHealthy := aggregator.WaitForHealthyState()
		if !aggregatorHealthy {
			logging.Debug("aggregator service not ready - restarting seems the only option")
			os.Exit(1)
		}

		// New context for aggregator streamer
		aggrContext, aggregatorCancel := context.WithCancel(context.Background())
		aristarchosStreamer, err := aggregator.CreateNewEntry(aggrContext)
		if err != nil {
			logging.Error(err.Error())
			return nil, nil, err
		}

		handler.Aggregator = aristarchosStreamer
		handler.AggregatorCancel = aggregatorCancel
	}

	return handler, conn, nil
}

func createEupalinosClient(serverAddress string) (pb.EupalinosClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewEupalinosClient(conn)
	return client, conn, nil
}
