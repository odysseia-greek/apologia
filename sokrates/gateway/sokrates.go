package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
)

type SokratesHandler struct {
	Cache      archytas.Client
	Streamer   pb.TraceServiceClient
	Cancel     context.CancelFunc
	Randomizer randomizer.Random
	MediaClient parti.
}
