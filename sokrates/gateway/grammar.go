package gateway

import (
	pbantisthenes "github.com/odysseia-greek/apologia/antisthenes/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/grammar"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateGrammarQuiz(request *pbantisthenes.CreationRequest, requestID, sessionId string) (*model.GrammarQuizResponse, error) {
	grammarClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()
	grpcResponse, err := s.GrammarClient.Question(grammarClientCtx, request)
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

func (s *SokratesHandler) CheckGrammar(request *pbantisthenes.AnswerRequest, requestID, sessionId string) (*model.GrammarAnswer, error) {
	grammarClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.GrammarClient.Answer(grammarClientCtx, request)
	if err != nil {
		return nil, err
	}

	return grammar.MapComprehensiveResponse(grpcResponse), nil
}

func (s *SokratesHandler) GrammarOptions(requestID, sessionId string) (*model.GrammarOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.GrammarClient.Options(optionsCtx, &pbantisthenes.OptionsRequest{})
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
