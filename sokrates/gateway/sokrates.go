package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/metadata"
	"time"
)

type SokratesHandler struct {
	Cache       archytas.Client
	Streamer    pbar.TraceService_ChorusClient
	Randomizer  randomizer.Random
	MediaClient *hedone.MediaClient
}

func (s *SokratesHandler) CreateMediaQuiz(theme, set, segment, quizType, order, requestID string, excludeWords []string) (*pbartrippos.QuizResponse, error) {
	request := pbartrippos.CreationRequest{
		Theme:        theme,
		Set:          set,
		Segment:      segment,
		QuizType:     quizType,
		Order:        order,
		ExcludeWords: excludeWords,
	}

	mediaClientCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestID})
	mediaClientCtx = metadata.NewOutgoingContext(context.Background(), md)

	return s.MediaClient.Question(mediaClientCtx, &request)

}
