package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/service"
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/multiplechoice"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateMultipleChoiceQuiz(request *pbkritias.CreationRequest, requestID string) (*model.MultipleChoiceResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MultiChoiceClient.Question(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &model.MultipleChoiceResponse{
		QuizItem:      &grpcResponse.QuizItem,
		NumberOfItems: &grpcResponse.NumberOfItems,
	}

	for _, opt := range grpcResponse.Options {
		quizResponse.Options = append(quizResponse.Options, &model.Options{
			Option: &opt.Option,
		})
	}

	for _, progress := range grpcResponse.Progress {
		quizResponse.Progress = append(quizResponse.Progress, &model.ProgressEntry{
			Greek:          &progress.Greek,
			Translation:    &progress.Translation,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckMultipleChoice(request *pbkritias.AnswerRequest, requestID string) (*model.ComprehensiveResponse, error) {
	multipleChoiceClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	multipleChoiceClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MultiChoiceClient.Answer(multipleChoiceClientCtx, request)
	if err != nil {
		return nil, err
	}

	return multiplechoice.MapComprehensiveResponse(grpcResponse), nil
}

func (s *SokratesHandler) MultipleChoiceOptions(requestID, sessionId string) (*model.AggregatedOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.MultiChoiceClient.Options(optionsCtx, &pbkritias.OptionsRequest{})
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
