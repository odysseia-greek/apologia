package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/anisthenes/kunismos"
	pb "github.com/odysseia-greek/apologia/anisthenes/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=Antisthenes
	logging.System(`
  ____  ____   ______  ____ _____ ______  __ __    ___  ____     ___  _____
 /    ||    \ |      ||    / ___/|      ||  |  |  /  _]|    \   /  _]/ ___/
|  o  ||  _  ||      | |  (   \_ |      ||  |  | /  [_ |  _  | /  [_(   \_ 
|     ||  |  ||_|  |_| |  |\__  ||_|  |_||  _  ||    _]|  |  ||    _]\__  |
|  _  ||  |  |  |  |   |  |/  \ |  |  |  |  |  ||   [_ |  |  ||   [_ /  \ |
|  |  ||  |  |  |  |   |  |\    |  |  |  |  |  ||     ||  |  ||     |\    |
|__|__||__|__|  |__|  |____|\___|  |__|  |__|__||_____||__|__||_____| \___|
`)
	logging.System("\"ἀρχὴ παιδεύσεως ἡ τῶν ὀνομάτων ἐπίσκεψις\"")
	logging.System("The investigation of the meaning of words is the beginning of education.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	config, err := kunismos.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(kunismos.GrammarInterceptor))

	pb.RegisterAnisthenesServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
