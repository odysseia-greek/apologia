package strategos

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/apologia/alkibiades/proto"
	"os"
	"time"
)

const (
	OPTIONSEGMENTKEY string = "archytassavedoptions"
)

func (j *JourneyServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := j.Elastic.Health().Info()
	dbHealth := &pb.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	return &pb.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
		Version:        os.Getenv("VERSION"),
	}, nil
}

func (j *JourneyServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := j.Archytas.Read(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()
		logging.Warn(fmt.Sprintf("%v", query))

		elasticResponse, err := j.Elastic.Query().MatchRaw(j.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = j.Archytas.Set(OPTIONSEGMENTKEY, string(elasticResponse))
		if err != nil {
			logging.Error(err.Error())
		}
	}

	return parseAggregationResult(unparsedResponse)
}
