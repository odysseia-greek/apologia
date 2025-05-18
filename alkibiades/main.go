package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/apologia/alkibiades/proto"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=alkibiades
	logging.System(`
  ____  _      __  _  ____  ____   ____   ____  ___      ___  _____
 /    || |    |  |/ ]|    ||    \ |    | /    ||   \    /  _]/ ___/
|  o  || |    |  ' /  |  | |  o  ) |  | |  o  ||    \  /  [_(   \_ 
|     || |___ |    \  |  | |     | |  | |     ||  D  ||    _]\__  |
|  _  ||     ||     \ |  | |  O  | |  | |  _  ||     ||   [_ /  \ |
|  |  ||     ||  .  | |  | |     | |  | |  |  ||     ||     |\    |
|__|__||_____||__|\_||____||_____||____||__|__||_____||_____| \___|

`)
	logging.System("\"ὦ Ἀλκιβιάδη, ἐπειδὴ περὶ τίνος Ἀθηναῖοι διανοοῦνται βουλεύεσθαι, ἀνίστασαι συμβουλεύσων;\"")
	logging.System("Alcibiades, on what subject do the Athenians propose to take advice, that you should stand up to advise them?")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	config, err := strategos.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(strategos.JourneyInterceptor))

	pb.RegisterAlkibiadesServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
