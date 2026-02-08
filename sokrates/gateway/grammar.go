package gateway

import (
	"context"

	v1 "github.com/odysseia-greek/apologia/antisthenes/gen/go/v1"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateGrammarQuiz(ctx context.Context, request *v1.CreationRequest) (*model.GrammarQuizResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.GrammarClient.CallWithReconnect(func(client *kunismos.GrammarClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	quizResponse := &model.GrammarQuizResponse{
		QuizItem:        &grpcResponse.QuizItem,
		NumberOfItems:   &grpcResponse.NumberOfItems,
		Stem:            &grpcResponse.Stem,
		DictionaryForm:  &grpcResponse.DictionaryForm,
		Translation:     &grpcResponse.Translation,
		Description:     &grpcResponse.Description,
		Difficulty:      &grpcResponse.Difficulty,
		ContractionRule: &grpcResponse.ContractionRule,
	}

	for _, opt := range grpcResponse.Options {
		quizResponse.Options = append(quizResponse.Options, &model.GrammarOption{
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

func (s *SokratesHandler) CheckGrammar(ctx context.Context, request *v1.AnswerRequest) (*model.GrammarAnswer, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.GrammarClient.CallWithReconnect(func(client *kunismos.GrammarClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	mappedResponse := &model.GrammarAnswer{
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

func (s *SokratesHandler) GrammarOptions(ctx context.Context) (*model.GrammarOptions, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.GrammarClient.CallWithReconnect(func(client *kunismos.GrammarClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(outCtx, &koinosv1.OptionsRequest{})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var themes []*model.GrammarThemes
	for _, grpcTheme := range grpcResponse.Themes {
		var segments []*model.GrammarSegment
		for _, grpcSegment := range grpcTheme.Segments {
			maxSet := int32(grpcSegment.MaxSet)
			segments = append(segments, &model.GrammarSegment{
				Name:       &grpcSegment.Name,
				MaxSet:     &maxSet,
				Difficulty: &grpcSegment.Difficulty,
			})
		}

		themes = append(themes, &model.GrammarThemes{
			Name:     &grpcTheme.Name,
			Segments: segments,
		})
	}

	return &model.GrammarOptions{
		Themes: themes,
	}, nil
}
