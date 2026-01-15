package gateway

import (
	"context"

	v1 "github.com/odysseia-greek/apologia/aristippos/gen/go/v1"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMediaQuiz(ctx context.Context, request *v1.CreationRequest) (*model.MediaQuizResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(outCtx, request)
		return innerErr
	})
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
			Translation:    &progress.Translation,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckMedia(ctx context.Context, request *v1.AnswerRequest) (*model.ComprehensiveResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	mappedResponse := &model.ComprehensiveResponse{
		Correct:  &grpcResponse.Correct,
		QuizWord: &grpcResponse.QuizWord,
		Finished: &grpcResponse.Finished,
	}

	for _, progress := range grpcResponse.Progress {
		mappedResponse.Progress = append(mappedResponse.Progress, &model.ProgressEntry{
			Greek:          &progress.Greek,
			Translation:    &progress.Translation,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return mappedResponse, nil
}

func (s *SokratesHandler) MediaOptions(ctx context.Context) (*model.AggregatedOptions, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(outCtx, &koinosv1.OptionsRequest{})
		return innerErr
	})
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
