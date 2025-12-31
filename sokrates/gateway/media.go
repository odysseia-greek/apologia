package gateway

import (
	v1 "github.com/odysseia-greek/apologia/aristippos/gen/go/v1"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateMediaQuiz(request *v1.CreationRequest, requestID, sessionId string) (*model.MediaQuizResponse, error) {
	mediaClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(mediaClientCtx, request)
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

func (s *SokratesHandler) CheckMedia(request *v1.AnswerRequest, requestID, sessionId string) (*model.ComprehensiveResponse, error) {
	mediaClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.AnswerResponse

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(mediaClientCtx, request)
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

func (s *SokratesHandler) MediaOptions(requestID, sessionId string) (*model.AggregatedOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.MediaClient.CallWithReconnect(func(client *hedone.MediaClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(optionsCtx, &v1.OptionsRequest{})
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
