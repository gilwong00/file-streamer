package grpctransport

import (
	"log"
	"net/http"

	"github.com/gilwong00/file-streamer/internal/gen/proto/v1/transferv1connect"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type connectRPCServer struct {
	transferv1connect.UnimplementedTransferServiceHandler
	address string
	storage *storage.BlobStorage
}

func NewConnectRPCServer(address string, bucketName string) (*connectRPCServer, error) {
	storage, err := storage.NewBlobStorage(bucketName)
	if err != nil {
		return nil, err
	}
	return &connectRPCServer{
		address: address,
		storage: storage,
	}, nil
}

func (s *connectRPCServer) StartServer() error {
	mux := http.NewServeMux()
	path, handler := transferv1connect.NewTransferServiceHandler(&connectRPCServer{})
	mux.Handle(path, handler)
	if err := http.ListenAndServe(
		s.address,
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatalf("listen failed: %v", err)
		return err
	}
	return nil
}
