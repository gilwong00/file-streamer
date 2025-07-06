package transport

import (
	"fmt"

	grpctransport "github.com/gilwong00/file-streamer/internal/server/transport/grpc"
	httptransport "github.com/gilwong00/file-streamer/internal/server/transport/http"
)

const (
	httpServerPort          = 3333
	connectRPCServerAddress = 6666
	fileDirectoryName       = "files"
)

func InitializeTransports() error {
	errors := make(chan error, 2)
	httpServer := httptransport.NewHttpServer(httpServerPort)
	connectRPCServer, err := grpctransport.NewConnectRPCServer(
		fmt.Sprintf(":%v", connectRPCServerAddress),
		fileDirectoryName,
	)
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
