package philia

import (
	"encoding/json"
	"fmt"
	pb "github.com/odysseia-greek/apologia/kriton/proto"
)

func quizAggregationQuery() map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"aggs": map[string]interface{}{
			"unique_themes": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "theme",
					"size":  1000,
				},
				"aggs": map[string]interface{}{
					"max_set": map[string]interface{}{
						"max": map[string]interface{}{
							"field": "set",
						},
					},
				},
			},
		},
	}
}

func parseAggregationResult(rawESOutput []byte) (*pb.AggregatedOptions, error) {
	// Define a structure to match the raw ES aggregation result format
	var esResponse struct {
		Aggregations struct {
			UniqueThemes struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
					MaxSet   struct {
						Value float64 `json:"value"`
					} `json:"max_set"`
				} `json:"buckets"`
			} `json:"unique_themes"`
		} `json:"aggregations"`
	}

	// Unmarshal the raw Elasticsearch output into the esResponse structure
	err := json.Unmarshal(rawESOutput, &esResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Elasticsearch response: %w", err)
	}

	var result pb.AggregatedOptions

	for _, themeBucket := range esResponse.Aggregations.UniqueThemes.Buckets {
		theme := &pb.Theme{
			Name:   themeBucket.Key,
			MaxSet: float32(themeBucket.MaxSet.Value),
		}
		result.Themes = append(result.Themes, theme)
	}

	return &result, nil
}
