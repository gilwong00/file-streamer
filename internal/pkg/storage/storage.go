package storage

import "context"

// Client defines the interface for interacting with an object storage service,
// supporting bucket creation and existence checks.
type Client interface {
	// CreateBucket attempts to create a new bucket with the given name.
	//
	// Returns ErrBucketAlreadyExists if the bucket already exists.
	// Returns any other error encountered during the creation process.
	CreateBucket(ctx context.Context, bucketName string) error

	// DoesBucketExists checks if a bucket with the specified name exists.
	//
	// Returns true if the bucket exists, false otherwise.
	// Returns an error if the existence check could not be performed.
	DoesBucketExists(ctx context.Context, bucketName string) (bool, error)
}

func NewStorageClient(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	useSSL bool,
) (Client, error) {
	return newClient(
		endpoint,
		accessKeyID,
		secretAccessKey,
		useSSL,
	)
}
