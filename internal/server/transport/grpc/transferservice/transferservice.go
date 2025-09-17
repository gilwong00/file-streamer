package transferservice

import (
	"github.com/gilwong00/file-streamer/internal/gen/proto/v1/transferv1connect"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
)

type transferService struct {
	storageClient storage.Client
	bucketName    string
}

func NewTransferService(
	storageClient storage.Client,
	bucketName string,
) transferv1connect.TransferServiceHandler {
	return &transferService{
		storageClient: storageClient,
		bucketName:    bucketName,
	}
}
