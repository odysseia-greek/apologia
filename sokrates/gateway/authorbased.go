package gateway

import (
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	pbxenofon "github.com/odysseia-greek/apologia/xenofon/proto"
)

func (s *SokratesHandler) CreateAuthorBasedQuiz(request *pbxenofon.CreationRequest, requestID, sessionId string) (*model.AuthorBasedResponse, error) {
	authorBasedCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.Question(authorBasedCtx, request)
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

func (s *SokratesHandler) CheckAuthorBased(request *pbxenofon.AnswerRequest, requestID, sessionId string) (*model.AuthorBasedAnswerResponse, error) {
	authorBasedCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.Answer(authorBasedCtx, request)
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

func (s *SokratesHandler) AuthorBasedOptions(requestID, sessionId string) (*model.AggregatedOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.Options(optionsCtx, &pbxenofon.OptionsRequest{})
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

func (s *SokratesHandler) AuthorBasedWordForms(request *pbxenofon.WordFormRequest, requestID, sessionId string) (*model.AuthorBasedWordFormsResponse, error) {
	wordFormsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.WordForms(wordFormsCtx, request)
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
