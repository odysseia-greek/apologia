package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/service"
	pbkriotn "github.com/odysseia-greek/apologia/kriton/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateDialogueQuiz(request *pbkriotn.CreationRequest, requestID string) (*model.DialogueQuizResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.DialogueClient.Question(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &model.DialogueQuizResponse{
		QuizMetadata: &model.QuizMetadata{
			Language: &grpcResponse.QuizMetadata.Language,
		},
		Theme:     &grpcResponse.Theme,
		Set:       &grpcResponse.Set,
		Segment:   &grpcResponse.Segment,
		Reference: &grpcResponse.Reference,
		Dialogue: &model.Dialogue{
			Introduction:  &grpcResponse.Dialogue.Introduction,
			Section:       &grpcResponse.Dialogue.Introduction,
			LinkToPerseus: &grpcResponse.Dialogue.Introduction,
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

func (s *SokratesHandler) CheckDialogueQuiz(request *pbkriotn.AnswerRequest, requestID string) (*model.DialogueAnswer, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.DialogueClient.Answer(mediaClientCtx, request)
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
