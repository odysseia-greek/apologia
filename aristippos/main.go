package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/apologia/aristippos/gen/go/v1"
	"github.com/odysseia-greek/apologia/aristippos/hedone"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ARISTIPPOS
	logging.System(`
  ____  ____   ____ _____ ______  ____  ____  ____   ___   _____
 /    ||    \ |    / ___/|      ||    ||    \|    \ /   \ / ___/
|  o  ||  D  ) |  (   \_ |      | |  | |  o  )  o  )     (   \_ 
|     ||    /  |  |\__  ||_|  |_| |  | |   _/|   _/|  O  |\__  |
|  _  ||    \  |  |/  \ |  |  |   |  | |  |  |  |  |     |/  \ |
|  |  ||  .  \ |  |\    |  |  |   |  | |  |  |  |  |     |\    |
|__|__||__|\_||____|\___|  |__|  |____||__|  |__|   \___/  \___|
`)
	logging.System("\"Κέκτημαι, οὐ κέκτημαι.\"")
	logging.System("I possess, I am not possessed")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := hedone.CreateNewConfig(ctx)
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

	v1.RegisterAristipposServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
