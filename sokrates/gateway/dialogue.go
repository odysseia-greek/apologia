package gateway

import (
	"github.com/odysseia-greek/apologia/kriton/philia"
	pbkriton "github.com/odysseia-greek/apologia/kriton/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) CreateDialogueQuiz(request *pbkriton.CreationRequest, requestID, sessionId string) (*model.DialogueQuizResponse, error) {
	dialogueClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *pbkriton.QuizResponse

	err := s.DialogueClient.CallWithReconnect(func(client *philia.DialogueClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(dialogueClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	quizResponse := &model.DialogueQuizResponse{
		QuizMetadata: &model.QuizMetadata{
			Language: &grpcResponse.QuizMetadata.Language,
		},
		Theme:     &grpcResponse.Theme,
		Set:       &grpcResponse.Set,
		Segment:   &grpcResponse.Dialogue.Section,
		Reference: &grpcResponse.Reference,
		Dialogue: &model.Dialogue{
			Introduction:  &grpcResponse.Dialogue.Introduction,
			Section:       &grpcResponse.Dialogue.Section,
			LinkToPerseus: &grpcResponse.Dialogue.LinkToPerseus,
		},
	}

	for _, speaker := range grpcResponse.Dialogue.Speakers {
		quizResponse.Dialogue.Speakers = append(quizResponse.Dialogue.Speakers, &model.Speaker{
			Name:        &speaker.Name,
			Shorthand:   &speaker.Shorthand,
			Translation: &speaker.Translation,
		})
	}

	for _, content := range grpcResponse.Content {
		quizResponse.Content = append(quizResponse.Content, &model.DialogueContent{
			Translation: &content.Translation,
			Greek:       &content.Greek,
			Place:       &content.Place,
			Speaker:     &content.Speaker,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckDialogueQuiz(request *pbkriton.AnswerRequest, requestID, sessionId string) (*model.DialogueAnswer, error) {
	dialogueClientCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *pbkriton.AnswerResponse

	err := s.DialogueClient.CallWithReconnect(func(client *philia.DialogueClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Answer(dialogueClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	answer := &model.DialogueAnswer{
		Percentage:    &grpcResponse.Percentage,
		Input:         nil,
		Answer:        nil,
		WronglyPlaced: nil,
	}

	for _, input := range grpcResponse.Input {
		answer.Input = append(answer.Input, &model.DialogueContent{
			Translation: &input.Translation,
			Greek:       &input.Greek,
			Place:       &input.Place,
			Speaker:     &input.Speaker,
		})
	}

	for _, output := range grpcResponse.Answer {
		answer.Answer = append(answer.Answer, &model.DialogueContent{
			Translation: &output.Translation,
			Greek:       &output.Greek,
			Place:       &output.Place,
			Speaker:     &output.Speaker,
		})
	}

	for _, placed := range grpcResponse.WronglyPlaced {
		answer.WronglyPlaced = append(answer.WronglyPlaced, &model.DialogueCorrection{
			Translation:  &placed.Translation,
			Greek:        &placed.Greek,
			Place:        &placed.Place,
			Speaker:      &placed.Speaker,
			CorrectPlace: &placed.CorrectPlace,
		})
	}

	return answer, nil
}

func (s *SokratesHandler) DialogueOptions(requestID, sessionId string) (*model.ThemedOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *pbkriton.AggregatedOptions

	err := s.DialogueClient.CallWithReconnect(func(client *philia.DialogueClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(optionsCtx, &pbkriton.OptionsRequest{})
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
