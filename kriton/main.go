package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/kriton/philia"
	pb "github.com/odysseia-greek/apologia/kriton/proto"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=KRITON
	logging.System(`
 __  _  ____   ____  ______   ___   ____  
|  |/ ]|    \ |    ||      | /   \ |    \ 
|  ' / |  D  ) |  | |      ||     ||  _  |
|    \ |    /  |  | |_|  |_||  O  ||  |  |
|     ||    \  |  |   |  |  |     ||  |  |
|  .  ||  .  \ |  |   |  |  |     ||  |  |
|__|\_||__|\_||____|  |__|   \___/ |__|__|
`)
	logging.System("\"ΣΩ. τί τηνικάδε ἀφῖξαι, ὦ Κρίτων; ἢ οὐ πρῲ ἔτι ἐστίν;\"")
	logging.System("Socrates: Why have you come at this time, Crito? Or isn’t it still early?")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	config, err := philia.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(philia.Interceptor))

	pb.RegisterKritonServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
