package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateMediaQuiz(request *pbartrippos.CreationRequest, requestID string) (*models.QuizResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MediaClient.Question(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &models.QuizResponse{
		QuizItem:      grpcResponse.QuizItem,
		NumberOfItems: int(grpcResponse.NumberOfItems),
	}

	for _, opt := range grpcResponse.Options {
		quizResponse.Options = append(quizResponse.Options, models.Options{
			Option:   opt.Option,
			AudioUrl: opt.AudioUrl,
			ImageUrl: opt.ImageUrl,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckMedia(request *pbartrippos.AnswerRequest, requestID string) (*models.ComprehensiveResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MediaClient.Answer(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	return mapComprehensiveResponse(grpcResponse), nil
}
