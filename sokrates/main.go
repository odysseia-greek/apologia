package main

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/apologia/sokrates/gateway"
	"log"
	"net/http"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/sokrates/routing"
)

const standardPort = ":8080"

func main() {
	// Set up port with environment variable fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	logging.System(`
  _____  ___   __  _  ____    ____  ______    ___  _____
 / ___/ /   \ |  |/ ]|    \  /    ||      |  /  _]/ ___/
(   \_ |     ||  ' / |  D  )|  o  ||      | /  [_(   \_ 
 \__  ||  O  ||    \ |    / |     ||_|  |_||    _]\__  |
 /  \ ||     ||     ||    \ |  _  |  |  |  |   [_ /  \ |
 \    ||     ||  .  ||  .  \|  |  |  |  |  |     |\    |
  \___| \___/ |__|\_||__|\_||__|__|  |__|  |_____| \___|
                                                        
`)
	logging.System("\"ἓν οἶδα ὅτι οὐδὲν οἶδα\"")
	logging.System("\"I know one thing, that I know nothing\"")
	logging.System("starting up and getting environment variables...")

	handler, err := gateway.CreateNewConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	graphqlServer := routing.InitRoutes(handler)

	logging.System(fmt.Sprintf("Server running on port %s", port))
	err = http.ListenAndServe(port, graphqlServer)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
