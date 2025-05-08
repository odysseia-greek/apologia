package gateway

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/antisthenes/kunismos"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type SokratesHandler struct {
	Streamer          pbar.TraceService_ChorusClient
	Randomizer        randomizer.Random
	MediaClient       *GenericGrpcClient[*hedone.MediaClient]
	MultiChoiceClient *GenericGrpcClient[*triakonta.MutpleChoiceClient]
	AuthorBasedClient *GenericGrpcClient[*anabasis.AuthorBasedClient]
	DialogueClient    *GenericGrpcClient[*philia.DialogueClient]
	GrammarClient     *GenericGrpcClient[*kunismos.GrammarClient]
	JourneyClient     *GenericGrpcClient[*strategos.JourneyClient]
}
type GenericGrpcClient[T any] struct {
	client  T
	address string
	dialFn  func(string) (T, error)
	mu      sync.Mutex
}

func (s *SokratesHandler) createRequestHeader(requestID, sessionId string) (context.Context, context.CancelFunc) {
	requestCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	md := metadata.New(map[string]string{config.HeaderKey: requestID,
		config.SessionIdKey: sessionId})
	requestCtx = metadata.NewOutgoingContext(requestCtx, md)

	return requestCtx, ctxCancel
}

func NewGenericGrpcClient[T any](address string, dialFn func(string) (T, error)) (*GenericGrpcClient[T], error) {
	client, err := dialFn(address)
	if err != nil {
		return nil, err
	}
	return &GenericGrpcClient[T]{
		client:  client,
		address: address,
		dialFn:  dialFn,
	}, nil
}

func (g *GenericGrpcClient[T]) Reconnect() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	client, err := g.dialFn(g.address)
	if err != nil {
		return err
	}
	g.client = client
	return nil
}

func (g *GenericGrpcClient[T]) CallWithReconnect(call func(client T) error) error {
	g.mu.Lock()
	client := g.client
	g.mu.Unlock()

	err := call(client)
	if err == nil {
		return nil
	}

	if !isConnectionError(err) {
		return err
	}

	// Log reconnecting event
	logging.Debug(fmt.Sprintf("connection error detected, reconnecting to %s", g.address))

	reconnectErr := g.Reconnect()
	if reconnectErr != nil {
		return reconnectErr
	}

	// Retry once
	g.mu.Lock()
	client = g.client
	g.mu.Unlock()

	return call(client)
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded
}
