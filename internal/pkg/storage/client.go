package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

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

// CreateBucket creates a new bucket with the specified name if it does not already exist.
//
// It first checks whether the bucket exists to avoid conflicts. If the bucket already exists,
// it returns ErrBucketAlreadyExists. Otherwise, it creates the bucket using MinIO client options.
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

// DoesBucketExists checks whether a bucket with the specified name exists.
//
// Returns true if the bucket exists, false otherwise. Returns an error
// if the existence check could not be performed.
func (b *blobStorageClient) DoesBucketExists(ctx context.Context, bucketName string) (bool, error) {
	return b.client.BucketExists(ctx, bucketName)
}

// GetObject retrieves the full object from the specified bucket.
//
// Returns a MinIO object (*minio.Object), allowing the caller to read the object
// or perform MinIO-specific operations. Returns an error if the object does not exist
// or cannot be accessed.
func (b *blobStorageClient) GetObject(
	ctx context.Context,
	bucketName,
	objectName string,
	opts minio.GetObjectOptions,
) (*minio.Object, error) {
	return b.client.GetObject(ctx, bucketName, objectName, opts)
}

// GetObjectWithRange retrieves a portion of the object using start and end byte offsets.
//
// This method is used for ranged streaming, resumable downloads, or partial reads.
// Returns an io.ReadCloser for the specified byte range. Returns an error if the range
// is invalid or the object cannot be accessed.
func (b *blobStorageClient) GetObjectWithRange(
	ctx context.Context,
	bucketName string,
	objectName string,
	start int64,
	end int64,
) (io.ReadCloser, error) {
	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}
	obj, err := b.client.GetObject(ctx, bucketName, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("getting ranged object: %w", err)
	}
	// Ensure object exists
	if _, err := obj.Stat(); err != nil {
		return nil, fmt.Errorf("stat object: %w", err)
	}
	return obj, nil
}

// GetObjectInfo retrieves metadata about the specified object.
//
// Returns an ObjectInfo containing the object's size and optionally other metadata fields.
// Returns an error if the object does not exist or cannot be accessed.
func (b *blobStorageClient) GetObjectInfo(
	ctx context.Context,
	bucketName string,
	objectName string,
) (ObjectInfo, error) {
	info, err := b.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return ObjectInfo{}, fmt.Errorf("stat object: %w", err)
	}
	return ObjectInfo{
		Size: info.Size,
	}, nil
}
