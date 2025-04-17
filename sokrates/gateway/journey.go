package gateway

import (
	palkibiades "github.com/odysseia-greek/apologia/alkibiades/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) JourneyOptions(requestID, sessionId string) (*model.JourneyOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.JourneyClient.Options(optionsCtx, &palkibiades.OptionsRequest{})
	if err != nil {
		return nil, err
	}

	var themes []*model.JourneyThemes
	for _, grpcTheme := range grpcResponse.Themes {
		var segments []*model.JourneySegment
		for _, grpcSegment := range grpcTheme.Segments {
			segments = append(segments, &model.JourneySegment{
				Name:     &grpcSegment.Name,
				Number:   &grpcSegment.Number,
				Location: &grpcSegment.Location,
				Coordinates: &model.Coordinates{
					X: float32ToFloat64Ptr(grpcSegment.Coordinates.X),
					Y: float32ToFloat64Ptr(grpcSegment.Coordinates.Y),
				},
			})
		}

		themes = append(themes, &model.JourneyThemes{
			Name:     &grpcTheme.Name,
			Segments: segments,
		})
	}

	return &model.JourneyOptions{
		Themes: themes,
	}, nil
}

func float32ToFloat64Ptr(f float32) *float64 {
	val := float64(f)
	return &val
}
