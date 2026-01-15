package gateway

import (
	"context"

	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/apologia/kritias/gen/go/v1"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMultipleChoiceQuiz(ctx context.Context, request *v1.CreationRequest) (*model.MultipleChoiceResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(outCtx, request)
		return innerErr
	})
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

func (s *SokratesHandler) CheckMultipleChoice(ctx context.Context, request *v1.AnswerRequest) (*model.ComprehensiveResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
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

func (s *SokratesHandler) MultipleChoiceOptions(ctx context.Context) (*model.ThemedOptions, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(outCtx, &koinosv1.OptionsRequest{})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var themes []*model.MultipleTheme
	for _, grpcTheme := range grpcResponse.Themes {
		maxSet := float64(grpcTheme.MaxSet)
		themes = append(themes, &model.MultipleTheme{
			Name:   &grpcTheme.Name,
			MaxSet: &maxSet,
		})
	}

	return &model.ThemedOptions{
		Themes: themes,
	}, nil
}
