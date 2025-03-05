package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/xenofon/anabasis"
	pb "github.com/odysseia-greek/apologia/xenofon/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const standardPort = ":50050"

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
	config, err := anabasis.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(anabasis.Interceptor))

	pb.RegisterXenofonServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
