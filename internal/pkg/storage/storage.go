package storage

type BlobStorage struct {
	bucketName string
}

func NewBlobStorage(bucketName string) (*BlobStorage, error) {
	return &BlobStorage{
		bucketName: bucketName,
	}, nil
}
