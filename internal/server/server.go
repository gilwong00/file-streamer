package server

import (
	"log"

	"github.com/gilwong00/file-streamer/internal/server/transport"
)

const (
	httpServerPort          = 3333
	connectRPCServerAddress = 6666
)

func StartServer() {
	if err := transport.InitializeTransports(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
