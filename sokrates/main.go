package main

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/apologia/sokrates/quiz"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
)

const standardPort = ":5000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=SOKRATES
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
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()

	sokratesConfig, err := quiz.CreateNewConfig(ctx)
	if err != nil {
		log.Print(err)
		log.Fatal("death has found me")
	}

	srv := quiz.InitRoutes(sokratesConfig)

	logging.System(fmt.Sprintf("%s : %s", "running on port", port))
	err = http.ListenAndServe(port, srv)
	if err != nil {
		panic(err)
	}
}
