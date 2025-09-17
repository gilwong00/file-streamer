package transport

import (
	"context"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
	grpctransport "github.com/gilwong00/file-streamer/internal/server/transport/grpc"
	httptransport "github.com/gilwong00/file-streamer/internal/server/transport/http"
)

// const (
// 	httpServerPort          = 3333
// 	connectRPCServerAddress = 6666
// 	fileDirectoryName       = "files"
// )

func InitializeTransports(
	ctx context.Context,
	config *config.Config,
	storageClient storage.Client,
) error {
	errors := make(chan error, 2)
	// TODO: add bucket name
	httpServer := httptransport.NewHttpServer(ctx, config, storageClient, "")
	connectRPCServer, err := grpctransport.NewConnectRPCServer(ctx, config, storageClient)
	if err != nil {
		return err
	}
	go func() {
		errors <- httpServer.Run()
	}()
	go func() {
		errors <- connectRPCServer.StartServer()
	}()
	// Wait for the first error
	err = <-errors
	return err
}
