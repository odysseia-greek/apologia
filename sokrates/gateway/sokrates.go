package gateway

import (
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/apologia/kriton/philia"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
)

type SokratesHandler struct {
	Streamer          pbar.TraceService_ChorusClient
	Randomizer        randomizer.Random
	MediaClient       *hedone.MediaClient
	MultiChoiceClient *triakonta.MutpleChoiceClient
	AuthorBasedClient *anabasis.AuthorBasedClient
	DialogueClient    *philia.DialogueClient
}
