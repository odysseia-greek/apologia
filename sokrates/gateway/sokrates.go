package gateway

import (
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
)

type SokratesHandler struct {
	Cache       archytas.Client
	Streamer    pbar.TraceService_ChorusClient
	Randomizer  randomizer.Random
	MediaClient *hedone.MediaClient
}
