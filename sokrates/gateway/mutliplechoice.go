package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/service"
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/multiplechoice"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateMultipleChoiceQuiz(request *pbkritias.CreationRequest, requestID string) (*model.MultipleChoiceResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

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

func (s *SokratesHandler) CheckMultipleChoice(request *pbkritias.AnswerRequest, requestID string) (*model.ComprehensiveResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MultiChoiceClient.Answer(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	return multiplechoice.MapComprehensiveResponse(grpcResponse), nil
}
