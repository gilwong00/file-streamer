package storage

type blobStorageClient struct {
	bucketName string
}

var _ Client = (*blobStorageClient)(nil)

func newClient(bucketName string) (*blobStorageClient, error) {
	// instanciate minio client
	return &blobStorageClient{
		bucketName: bucketName,
	}, nil
}
