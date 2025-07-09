package server

import (
	"context"
	"log"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/server/transport"
)

func StartServer(ctx context.Context, config *config.Config) error {
	if err := transport.InitializeTransports(ctx, config); err != nil {
		log.Printf("server error: %v", err)
		return err
	}
	return nil
}
