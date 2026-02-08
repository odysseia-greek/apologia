package gateway

import (
	"context"

	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	v1 "github.com/odysseia-greek/apologia/xenofon/gen/go/v1"
)

func (s *SokratesHandler) CreateAuthorBasedQuiz(ctx context.Context, request *v1.CreationRequest) (*model.AuthorBasedResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.AuthorBasedClient.CallWithReconnect(func(client *anabasis.AuthorBasedClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	quizResponse := &model.AuthorBasedResponse{
		FullSentence: &grpcResponse.FullSentence,
		Translation:  &grpcResponse.Translation,
		Reference:    &grpcResponse.Reference,
		Quiz: &model.AuthorBasedQuiz{
			QuizItem:      &grpcResponse.Quiz.QuizItem,
			NumberOfItems: &grpcResponse.Quiz.NumberOfItems,
		},
	}

	for _, option := range grpcResponse.Quiz.Options {
		quizResponse.Quiz.Options = append(quizResponse.Quiz.Options, &model.AuthorBasedOptions{QuizWord: &option.QuizWord})
	}

	for _, opt := range grpcResponse.GrammarQuiz {
		grammarQuiz := &model.GrammarQuizAdded{
			CorrectAnswer:    &opt.CorrectAnswer,
			WordInText:       &opt.WordInText,
			ExtraInformation: &opt.ExtraInformation,
		}

		for _, option := range opt.Options {
			grammarQuiz.Options = append(grammarQuiz.Options, &model.AuthorBasedOptions{QuizWord: &option.QuizWord})
		}

		quizResponse.GrammarQuiz = append(quizResponse.GrammarQuiz, grammarQuiz)
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

func (s *SokratesHandler) CheckAuthorBased(
	ctx context.Context,
	request *v1.AnswerRequest,
) (*model.AuthorBasedAnswerResponse, error) {

	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.AuthorBasedClient.CallWithReconnect(func(client *anabasis.AuthorBasedClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	answerResponse := &model.AuthorBasedAnswerResponse{
		Correct:     &grpcResponse.Correct,
		QuizWord:    &grpcResponse.QuizWord,
		Finished:    &grpcResponse.Finished,
		WordsInText: convertStringSliceToPointer(grpcResponse.WordsInText),
	}

	for _, progress := range grpcResponse.Progress {
		answerResponse.Progress = append(answerResponse.Progress, &model.ProgressEntry{
			Greek:          &progress.Greek,
			Translation:    &progress.Translation,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return answerResponse, nil
}

func (s *SokratesHandler) AuthorBasedOptions(ctx context.Context) (*model.AggregatedOptions, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.AuthorBasedClient.CallWithReconnect(func(client *anabasis.AuthorBasedClient) error {
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

func (s *SokratesHandler) AuthorBasedWordForms(ctx context.Context, request *v1.WordFormRequest) (*model.AuthorBasedWordFormsResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.WordFormResponse

	err := s.AuthorBasedClient.CallWithReconnect(func(client *anabasis.AuthorBasedClient) error {
		var innerErr error
		grpcResponse, innerErr = client.WordForms(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	forms := &model.AuthorBasedWordFormsResponse{
		Forms: make([]*model.AuthorBasedWordForm, len(grpcResponse.Forms)),
	}
	var wordForms []*model.AuthorBasedWordForm
	for _, grpcForm := range grpcResponse.Forms {
		wordForm := &model.AuthorBasedWordForm{
			DictionaryForm: &grpcForm.DictionaryForm,
			WordsInText:    convertStringSliceToPointer(grpcForm.WordsInText),
		}
		wordForms = append(wordForms, wordForm)
	}

	forms.Forms = wordForms

	return forms, nil
}

func convertStringSliceToPointer(strings []string) []*string {
	var ptrSlice []*string
	for _, s := range strings {
		ptrSlice = append(ptrSlice, &s)
	}
	return ptrSlice
}
