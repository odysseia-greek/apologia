package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/apologia/alkibiades/gen/go/v1"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc"
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
	cfg, err := strategos.CreateNewConfig(ctx)
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

	v1.RegisterAlkibiadesServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
