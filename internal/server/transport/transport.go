package transport

import (
	"context"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	grpctransport "github.com/gilwong00/file-streamer/internal/server/transport/grpc"
	httptransport "github.com/gilwong00/file-streamer/internal/server/transport/http"
)

const (
	httpServerPort          = 3333
	connectRPCServerAddress = 6666
	fileDirectoryName       = "files"
)

func InitializeTransports(ctx context.Context, config *config.Config) error {
	errors := make(chan error, 2)
	httpServer := httptransport.NewHttpServer(ctx, config)
	connectRPCServer, err := grpctransport.NewConnectRPCServer(ctx, config)
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
