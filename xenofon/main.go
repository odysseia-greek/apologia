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

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=XENOFON
	logging.System(`
 __ __    ___  ____    ___   _____   ___   ____  
|  |  |  /  _]|    \  /   \ |     | /   \ |    \ 
|  |  | /  [_ |  _  ||     ||   __||     ||  _  |
|_   _||    _]|  |  ||  O  ||  |_  |  O  ||  |  |
|     ||   [_ |  |  ||     ||   _] |     ||  |  |
|  |  ||     ||  |  ||     ||  |   |     ||  |  |
|__|__||_____||__|__| \___/ |__|    \___/ |__|__|
`)
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

}
