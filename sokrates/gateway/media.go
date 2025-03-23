package gateway

import (
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/media"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMediaQuiz(request *pbartrippos.CreationRequest, requestID, sessionId string) (*model.MediaQuizResponse, error) {
	mediaClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()
	grpcResponse, err := s.MediaClient.Question(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &model.MediaQuizResponse{
		QuizItem:      &grpcResponse.QuizItem,
		NumberOfItems: &grpcResponse.NumberOfItems,
	}

	for _, opt := range grpcResponse.Options {
		quizResponse.Options = append(quizResponse.Options, &model.MediaOptions{
			Option:   &opt.Option,
			AudioURL: &opt.AudioUrl,
			ImageURL: &opt.ImageUrl,
		})
	}

	for _, progress := range grpcResponse.Progress {
		quizResponse.Progress = append(quizResponse.Progress, &model.ProgressEntry{
			Greek:          &progress.Greek,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckMedia(request *pbartrippos.AnswerRequest, requestID, sessionId string) (*model.ComprehensiveResponse, error) {
	mediaClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.MediaClient.Answer(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	return media.MapComprehensiveResponse(grpcResponse), nil
}

func (s *SokratesHandler) MediaOptions(requestID, sessionId string) (*model.AggregatedOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.MediaClient.Options(optionsCtx, &pbartrippos.OptionsRequest{})
	if err != nil {
		return nil, err
	}

	var themes []*model.Theme
	for _, grpcTheme := range grpcResponse.Themes {
		var segments []*model.Segment
		for _, grpcSegment := range grpcTheme.Segments {
			maxSet := float64(grpcSegment.MaxSet)
			segments = append(segments, &model.Segment{
				Name:   &grpcSegment.Name,
				MaxSet: &maxSet,
			})
		}

		themes = append(themes, &model.Theme{
			Name:     &grpcTheme.Name,
			Segments: segments,
		})
	}

	return &model.AggregatedOptions{
		Themes: themes,
	}, nil
}
