package transport

import (
	"fmt"

	grpctransport "github.com/gilwong00/file-streamer/internal/server/transport/grpc"
	httptransport "github.com/gilwong00/file-streamer/internal/server/transport/http"
)

const (
	httpServerPort          = 3333
	connectRPCServerAddress = 6666
)

func InitializeTransports() error {
	errs := make(chan error, 2)
	httpServer := httptransport.NewHttpServer(httpServerPort)
	connectRPCServer, err := grpctransport.NewConnectRPCServer(fmt.Sprintf(":%v", connectRPCServerAddress), "files")
	if err != nil {
		return err
	}

	go func() {
		errs <- httpServer.Run()
	}()

	go func() {
		errs <- connectRPCServer.StartServer()
	}()

	// Wait for the first error
	err = <-errs
	return err
}
