package kunismos

import (
	"encoding/json"
	"fmt"
	pb "github.com/odysseia-greek/apologia/anisthenes/proto"
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
					"unique_segments": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "segment",
							"size":  1000,
						},
						"aggs": map[string]interface{}{
							"max_set": map[string]interface{}{
								"max": map[string]interface{}{
									"field": "set",
								},
							},
							"difficulty": map[string]interface{}{
								"top_hits": map[string]interface{}{
									"size": 1,
									"_source": map[string]interface{}{
										"includes": []string{"difficulty"},
									},
								},
							},
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
					Key            string `json:"key"`
					DocCount       int    `json:"doc_count"`
					UniqueSegments struct {
						Buckets []struct {
							Key      string `json:"key"`
							DocCount int    `json:"doc_count"`
							MaxSet   struct {
								Value float64 `json:"value"`
							} `json:"max_set"`
							Difficulty struct {
								Hits struct {
									Hits []struct {
										Source struct {
											Difficulty string `json:"difficulty"`
										} `json:"_source"`
									} `json:"hits"`
								} `json:"hits"`
							} `json:"difficulty"`
						} `json:"buckets"`
					} `json:"unique_segments"`
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
			Name: themeBucket.Key,
		}
		for _, segmentBucket := range themeBucket.UniqueSegments.Buckets {
			// Safely pull difficulty if available
			difficulty := ""
			if len(segmentBucket.Difficulty.Hits.Hits) > 0 {
				difficulty = segmentBucket.Difficulty.Hits.Hits[0].Source.Difficulty
			}

			segment := &pb.Segment{
				Name:       segmentBucket.Key,
				Difficulty: difficulty,
				MaxSet:     float32(segmentBucket.MaxSet.Value),
			}
			theme.Segments = append(theme.Segments, segment)
		}
		result.Themes = append(result.Themes, theme)
	}

	return &result, nil
}
