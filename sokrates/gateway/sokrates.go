package gateway

import (
	"context"
	"time"

	"github.com/odysseia-greek/agora/hesiodos"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/aspasia/rhetorike"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc/metadata"
)

type SokratesHandler struct {
	Streamer          arv1.TraceService_ChorusClient
	Randomizer        randomizer.Random
	MediaClient       *hesiodos.GenericGrpcClient[*hedone.MediaClient]
	MultiChoiceClient *hesiodos.GenericGrpcClient[*triakonta.MutpleChoiceClient]
	AuthorBasedClient *hesiodos.GenericGrpcClient[*anabasis.AuthorBasedClient]
	DialogueClient    *hesiodos.GenericGrpcClient[*philia.DialogueClient]
	GrammarClient     *hesiodos.GenericGrpcClient[*kunismos.GrammarClient]
	JourneyClient     *hesiodos.GenericGrpcClient[*strategos.JourneyClient]
	GathererClient    *hesiodos.GenericGrpcClient[*rhetorike.GathererClient]
}

func (s *SokratesHandler) outgoingCtx(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)

	reqID, _ := parent.Value(config.HeaderKey).(string)
	sessionID, _ := parent.Value(config.SessionIdKey).(string)

	kvs := make([]string, 0, 4)

	if reqID != "" {
		kvs = append(kvs, config.HeaderKey, reqID)
	}
	if sessionID != "" {
		kvs = append(kvs, config.SessionIdKey, sessionID)
	}

	if len(kvs) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, kvs...)
	}

	return ctx, cancel
}
