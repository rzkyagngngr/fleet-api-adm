package helper

import (
	"context"
	"time"
)

// StorageProvider defines the interface for interacting with S3-like storage
type StorageProvider interface {
	GeneratePresignedPutURL(ctx context.Context, bucket, key string, contentType string, ttl time.Duration) (string, error)
	GeneratePresignedGetURL(ctx context.Context, bucket, key string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, bucket, key string) error
}
