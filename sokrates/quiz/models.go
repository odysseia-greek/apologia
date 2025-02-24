package quiz

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
)

func quizAggregationQuery(quizType string) map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"quizType": quizType,
			},
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
						},
					},
				},
			},
		},
	}
}

func parseAggregationResult(rawESOutput []byte, quizType string) (*models.AggregatedOptions, error) {
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

	var result models.AggregatedOptions

	switch quizType {
	case models.MEDIA:
		for _, themeBucket := range esResponse.Aggregations.UniqueThemes.Buckets {
			theme := models.Theme{
				Name: themeBucket.Key,
			}
			for _, segmentBucket := range themeBucket.UniqueSegments.Buckets {
				segment := models.Segment{
					Name:   segmentBucket.Key,
					MaxSet: segmentBucket.MaxSet.Value,
				}
				theme.Segments = append(theme.Segments, segment)
			}
			result.Themes = append(result.Themes, theme)
		}
	case models.AUTHORBASED:
		for _, themeBucket := range esResponse.Aggregations.UniqueThemes.Buckets {
			theme := models.Theme{
				Name: themeBucket.Key,
			}
			for _, segmentBucket := range themeBucket.UniqueSegments.Buckets {
				segment := models.Segment{
					Name:   segmentBucket.Key,
					MaxSet: segmentBucket.MaxSet.Value,
				}
				theme.Segments = append(theme.Segments, segment)
			}
			result.Themes = append(result.Themes, theme)
		}
	case models.DIALOGUE:
		for _, themeBucket := range esResponse.Aggregations.UniqueThemes.Buckets {
			theme := models.Theme{
				Name: themeBucket.Key,
			}
			segment := models.Segment{
				Name:   themeBucket.Key,
				MaxSet: float64(themeBucket.DocCount),
			}
			theme.Segments = append(theme.Segments, segment)
			result.Themes = append(result.Themes, theme)
		}

	case models.MULTICHOICE:
		for _, themeBucket := range esResponse.Aggregations.UniqueThemes.Buckets {
			theme := models.Theme{
				Name: themeBucket.Key,
			}
			segment := models.Segment{
				Name:   themeBucket.Key,
				MaxSet: float64(themeBucket.DocCount),
			}
			theme.Segments = append(theme.Segments, segment)
			result.Themes = append(result.Themes, theme)
		}
	}

	return &result, nil
}
