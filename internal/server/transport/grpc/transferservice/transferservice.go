package transferservice

import (
	"github.com/gilwong00/file-streamer/internal/gen/proto/v1/transferv1connect"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
)

type transferService struct {
	storageClient *storage.Client
}

func NewTransferService(storageClient *storage.Client) transferv1connect.TransferServiceHandler {
	return &transferService{
		storageClient: storageClient,
	}
}
