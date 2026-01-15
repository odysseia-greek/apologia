package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/apologia/kritias/gen/go/v1"
	"github.com/odysseia-greek/apologia/kritias/triakonta"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=KRITIAS
	logging.System(`
 __  _  ____   ____  ______  ____   ____  _____
|  |/ ]|    \ |    ||      ||    | /    |/ ___/
|  ' / |  D  ) |  | |      | |  | |  o  (   \_ 
|    \ |    /  |  | |_|  |_| |  | |     |\__  |
|     ||    \  |  |   |  |   |  | |  _  |/  \ |
|  .  ||  .  \ |  |   |  |   |  | |  |  |\    |
|__|\_||__|\_||____|  |__|  |____||__|__| \___|
`)
	logging.System("\"ὅτι ὑικὸν αὐτῷ δοκοίη πάσχειν ὁ Κριτίας, ἐπιθυμῶν Εὐθυδήμῳ προσκνῆσθαι\nὥσπερ τὰ ὕδια τοῖς λίθοις.\"")
	logging.System("Critias seems to have the feelings of a pig: he can no more keep away from Euthydemus than pigs can help rubbing themselves against stones.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := triakonta.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(
		grpc.UnaryInterceptor(
			comedy.UnaryServerInterceptor(
				cfg.Streamer,
				comedy.WithHeaderKey(config.HeaderKey),
				comedy.WithContextKeyName(config.DefaultTracingName),
				comedy.WithCloseHop(),
			),
		),
	)

	v1.RegisterKritiasServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
