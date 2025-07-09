package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// blobStorageClient is a MinIO-backed implementation of the Client interface.
type blobStorageClient struct {
	client *minio.Client
}

// Compile-time check to ensure blobStorageClient implements Client.
var _ Client = (*blobStorageClient)(nil)

// ErrBucketAlreadyExists is returned when CreateBucket is called on an existing bucket.
var ErrBucketAlreadyExists = errors.New("bucket already exists")

// newClient initializes and returns a new blobStorageClient configured to connect
// to the object storage endpoint using the provided credentials and SSL setting.
func newClient(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	useSSL bool,
) (*blobStorageClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &blobStorageClient{
		client: client,
	}, nil
}

// CreateBucket creates a bucket with the given name if it does not already exist.
//
// It first checks if the bucket exists to avoid conflicts.
func (b *blobStorageClient) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := b.DoesBucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("checking if bucket exists: %w", err)
	}
	if exists {
		return ErrBucketAlreadyExists
	}
	return b.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

// DoesBucketExists checks whether the specified bucket exists.
func (b *blobStorageClient) DoesBucketExists(ctx context.Context, bucketName string) (bool, error) {
	return b.client.BucketExists(ctx, bucketName)
}
