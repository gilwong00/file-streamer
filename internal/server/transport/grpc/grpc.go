package grpctransport

import (
	"log"
	"net/http"

	"github.com/gilwong00/file-streamer/internal/gen/proto/v1/transferv1connect"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
	"github.com/gilwong00/file-streamer/internal/server/transport/grpc/transferservice"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type connectRPCServer struct {
	transferv1connect.UnimplementedTransferServiceHandler
	address       string
	storageClient *storage.Client
}

func NewConnectRPCServer(address string, bucketName string) (*connectRPCServer, error) {
	storageClient, err := storage.NewStorageClient(bucketName)
	if err != nil {
		return nil, err
	}
	return &connectRPCServer{
		address:       address,
		storageClient: &storageClient,
	}, nil
}

func (s *connectRPCServer) StartServer() error {
	mux := http.NewServeMux()
	transferService := transferservice.NewTransferService(s.storageClient)
	transferPath, transferHandler := transferv1connect.NewTransferServiceHandler(transferService)
	mux.Handle(transferPath, transferHandler)

	log.Printf("ConnectRPC server listening on %s", s.address)
	if err := http.ListenAndServe(
		s.address,
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatalf("listen failed: %v", err)
		return err
	}
	log.Printf("ConnectRPC server stopped")
	return nil
}
