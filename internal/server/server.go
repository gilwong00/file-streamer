package server

import (
	"context"
	"log"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
	"github.com/gilwong00/file-streamer/internal/server/transport"
)

func StartServer(ctx context.Context, config *config.Config) error {
	storageClient, err := storage.NewStorageClient(
		config.MinioHost,
		config.MinioAccessKeyID,
		config.MinioAccessKey,
		config.MinioUseSSL,
	)
	if err != nil {
		return err
	}
	if err := transport.InitializeTransports(ctx, config, storageClient); err != nil {
		log.Printf("server error: %v", err)
		return err
	}
	return nil
}
