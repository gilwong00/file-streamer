package server

import (
	"log"

	"github.com/gilwong00/file-streamer/internal/server/transport"
)

func StartServer() error {
	if err := transport.InitializeTransports(); err != nil {
		log.Printf("server error: %v", err)
		return err
	}
	return nil
}
