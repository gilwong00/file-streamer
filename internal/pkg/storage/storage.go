package storage

type Client interface {
}

func NewStorageClient(bucketName string) (Client, error) {
	return newClient(bucketName)
}
