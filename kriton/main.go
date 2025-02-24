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
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

}
