package gateway

import (
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/multiplechoice"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMultipleChoiceQuiz(request *pbkritias.CreationRequest, requestID, sessionId string) (*model.MultipleChoiceResponse, error) {
	mediaClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

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

func (s *SokratesHandler) CheckMultipleChoice(request *pbkritias.AnswerRequest, requestID, sessionId string) (*model.ComprehensiveResponse, error) {
	multipleChoiceClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.MultiChoiceClient.Answer(multipleChoiceClientCtx, request)
	if err != nil {
		return nil, err
	}

	return multiplechoice.MapComprehensiveResponse(grpcResponse), nil
}

func (s *SokratesHandler) MultipleChoiceOptions(requestID, sessionId string) (*model.MultipleChoiceOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.MultiChoiceClient.Options(optionsCtx, &pbkritias.OptionsRequest{})
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

	return &model.MultipleChoiceOptions{
		Themes: themes,
	}, nil
}
