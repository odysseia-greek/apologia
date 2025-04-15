package strategos

import (
	"encoding/json"
	"fmt"
	pb "github.com/odysseia-greek/apologia/alkibiades/proto"
)

func quizAggregationQuery() map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
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
							"segment_data": map[string]interface{}{
								"top_hits": map[string]interface{}{
									"size": 1,
									"_source": map[string]interface{}{
										"includes": []string{
											"number",
											"location",
											"coordinates",
										},
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
	// Define a structure to match the new Elasticsearch response
	var esResponse struct {
		Aggregations struct {
			UniqueThemes struct {
				Buckets []struct {
					Key            string `json:"key"` // theme
					DocCount       int    `json:"doc_count"`
					UniqueSegments struct {
						Buckets []struct {
							Key         string `json:"key"` // segment
							DocCount    int    `json:"doc_count"`
							SegmentData struct {
								Hits struct {
									Hits []struct {
										Source struct {
											Number      int32  `json:"number"`
											Location    string `json:"location"`
											Coordinates struct {
												X float32 `json:"x"`
												Y float32 `json:"y"`
											} `json:"coordinates"`
										} `json:"_source"`
									} `json:"hits"`
								} `json:"hits"`
							} `json:"segment_data"`
						} `json:"buckets"`
					} `json:"unique_segments"`
				} `json:"buckets"`
			} `json:"unique_themes"`
		} `json:"aggregations"`
	}

	// Unmarshal Elasticsearch response
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
			var coords *pb.Coordinates
			if len(segmentBucket.SegmentData.Hits.Hits) > 0 {
				src := segmentBucket.SegmentData.Hits.Hits[0].Source
				coords = &pb.Coordinates{
					X: src.Coordinates.X,
					Y: src.Coordinates.Y,
				}

				segment := &pb.Segments{
					Name:        segmentBucket.Key,
					Number:      src.Number,
					Location:    src.Location,
					Coordinates: coords,
				}
				theme.Segments = append(theme.Segments, segment)
			}
		}

		result.Themes = append(result.Themes, theme)
	}

	return &result, nil
}
