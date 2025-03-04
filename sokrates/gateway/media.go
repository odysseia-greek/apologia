package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/service"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateMediaQuiz(request *pbartrippos.CreationRequest, requestID string) (*model.QuizResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	grpcResponse, err := s.MediaClient.Question(mediaClientCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &model.QuizResponse{
		QuizItem:      &grpcResponse.QuizItem,
		NumberOfItems: &grpcResponse.NumberOfItems,
	}

	for _, opt := range grpcResponse.Options {
		quizResponse.Options = append(quizResponse.Options, &model.Options{
			Option:   &opt.Option,
			AudioURL: &opt.AudioUrl,
			ImageURL: &opt.ImageUrl,
		})
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckMedia(request *pbartrippos.AnswerRequest, requestID string) (*model.ComprehensiveResponse, error) {
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
