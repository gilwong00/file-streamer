package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// ObjectInfo contains metadata about an object in storage.
type ObjectInfo struct {
	Size int64
	// You can add other metadata fields like ContentType, ETag, LastModified, etc.
}

// GetObjectOptions defines optional parameters for retrieving an object
// from object storage.
//
// Start and End specify a byte range to fetch a portion of the object.
// If both are zero (the default), the entire object will be retrieved.
// End is inclusive, so the range includes bytes from Start to End.
type GetObjectOptions struct {
	Start int64 // Starting byte offset (inclusive)
	End   int64 // Ending byte offset (inclusive)
}

// Client defines the interface for interacting with an object storage service,
// supporting bucket operations and object retrieval.
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

	// GetObject retrieves the full object from the bucket.
	//
	// Returns a MinIO object (*minio.Object) to allow reading and other MinIO-specific operations.
	// Returns an error if the object does not exist or cannot be accessed.
	GetObject(ctx context.Context, bucketName, objectName string, opts GetObjectOptions) (*minio.Object, error)

	// GetObjectWithRange retrieves a portion of the object using start and end byte offsets.
	//
	// This is used for resumable downloads or partial reads.
	// Returns an io.ReadCloser for the specified byte range.
	// Returns an error if the range is invalid or the object cannot be accessed.
	GetObjectWithRange(ctx context.Context, bucketName, objectName string, start, end int64) (io.ReadCloser, error)

	// GetObjectInfo retrieves metadata about the specified object.
	//
	// Returns an ObjectInfo containing the object's size and other optional metadata.
	// Returns an error if the object does not exist or cannot be accessed.
	GetObjectInfo(ctx context.Context, bucketName, objectName string) (ObjectInfo, error)
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
