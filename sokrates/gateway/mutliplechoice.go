package gateway

import (
	v1 "github.com/odysseia-greek/apologia/kritias/gen/go/v1"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMultipleChoiceQuiz(request *v1.CreationRequest, requestID, sessionId string) (*model.MultipleChoiceResponse, error) {
	multipleChoiceCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(multipleChoiceCtx, request)
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

func (s *SokratesHandler) CheckMultipleChoice(request *v1.AnswerRequest, requestID, sessionId string) (*model.ComprehensiveResponse, error) {
	multipleChoiceCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(multipleChoiceCtx, request)
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

	return mappedResponse, nil
}

func (s *SokratesHandler) MultipleChoiceOptions(requestID, sessionId string) (*model.ThemedOptions, error) {
	multipleChoiceCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.MultiChoiceClient.CallWithReconnect(func(client *triakonta.MutpleChoiceClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(multipleChoiceCtx, &v1.OptionsRequest{})
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
