package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	"github.com/odysseia-greek/apologia/aspasia/rhetorike"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ASPASIA&x=none&v=4&h=4&w=80&we=false
	logging.System(`
  ____  _____ ____   ____  _____ ____   ____ 
 /    |/ ___/|    \ /    |/ ___/|    | /    |
|  o  (   \_ |  o  )  o  (   \_  |  | |  o  |
|     |\__  ||   _/|     |\__  | |  | |     |
|  _  |/  \ ||  |  |  _  |/  \ | |  | |  _  |
|  |  |\    ||  |  |  |  |\    | |  | |  |  |
|__|__| \___||__|  |__|__| \___||____||__|__|
`)
	logging.System("\"ἐπεὶ δʼ Ἀσπασία χαριζόμενος δοκεῖ πρᾶξαι τὰ πρὸς Σαμίους, ἐνταῦθα ἂν εἴη καιρὸς διαπορῆσαι μάλιστα περὶ τῆς ἀνθρώπου, τίνα τέχνην ἢ δύναμιν τοσαύτην ἔχουσα τῶν τε πολιτικῶν τοὺς πρωτεύοντας ἐχειρώσατο καὶ τοῖς φιλοσόφοις οὐ φαῦλον οὐδʼ ὀλίγον ὑπὲρ αὑτῆς παρέσχε λόγον.\"")
	logging.System("Now, since it is thought that he proceeded thus against the Samians to gratify Aspasia, this may be a fitting place to raise the query what great art or power this woman had, that she managed as she pleased the foremost men of the state, and afforded the philosophers occasion to discuss her in exalted terms and at great length.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := rhetorike.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}
	defer func() {
		if err := cfg.Close(); err != nil {
			logging.Error(fmt.Sprintf("failed to close Aspasia resources: %v", err))
		}
	}()

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

	reflection.Register(server)

	v1.RegisterAspasiaServiceServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
