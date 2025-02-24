package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
	"os"
)

const standardPort = ":50050"

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
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

}
