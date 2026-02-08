package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	v1 "github.com/odysseia-greek/apologia/xenofon/gen/go/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=XENOFON
	logging.System(`
 __ __    ___  ____    ___   _____   ___   ____  
|  |  |  /  _]|    \  /   \ |     | /   \ |    \ 
|  |  | /  [_ |  _  ||     ||   __||     ||  _  |
|_   _||    _]|  |  ||  O  ||  |_  |  O  ||  |  |
|     ||   [_ |  |  ||     ||   _] |     ||  |  |
|  |  ||     ||  |  ||     ||  |   |     ||  |  |
|__|__||_____||__|__| \___/ |__|    \___/ |__|__|
`)
	logging.System("\"ἀλλὰ μὴν καλόν τε καὶ δίκαιον καὶ ὅσιον καὶ ἥδιον τῶν ἀγαθῶν μᾶλλον ἢ τῶν κακῶν μεμνῆσθαι.\"")
	logging.System("Yet surely it is more honourable and fair, more righteous and gracious to remember good deeds than evil.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := anabasis.CreateNewConfig(ctx)
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

	v1.RegisterXenofonServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
