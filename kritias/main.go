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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=KRITIAS
	logging.System(`
 __  _  ____   ____  ______  ____   ____  _____
|  |/ ]|    \ |    ||      ||    | /    |/ ___/
|  ' / |  D  ) |  | |      | |  | |  o  (   \_ 
|    \ |    /  |  | |_|  |_| |  | |     |\__  |
|     ||    \  |  |   |  |   |  | |  _  |/  \ |
|  .  ||  .  \ |  |   |  |   |  | |  |  |\    |
|__|\_||__|\_||____|  |__|  |____||__|__| \___|
`)
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

}
