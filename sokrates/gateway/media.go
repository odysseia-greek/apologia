package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/plato/service"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	"github.com/odysseia-greek/apologia/sokrates/gateway/media"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *SokratesHandler) CreateMediaQuiz(request *pbartrippos.CreationRequest, requestID, sessionId string) (*model.MediaQuizResponse, error) {
	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID,
		"session-id": sessionId})
	mediaClientCtx = metadata.NewOutgoingContext(mediaClientCtx, md)
	
	grpcResponse, err := s.MediaClient.Question(mediaClientCtx, request)
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

	return media.MapComprehensiveResponse(grpcResponse), nil
}
