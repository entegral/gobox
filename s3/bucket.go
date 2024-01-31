package s3

import (
	"context"

	"github.com/entegral/gobox/types"
)

// BucketManager is a struct that can be embedded into other structs to provide s3 functionality
type BucketManager struct {
	Bucket string
}

// NewBucketManager returns a new BucketManager with the provided bucket
func NewBucketManager(bucket string) *BucketManager {
	return &BucketManager{
		Bucket: bucket,
	}
}

func (b *BucketManager) PutObject(ctx context.Context, item types.Keyable) error {
	err := Put(ctx, b.Bucket, item)
	return err
}

func (b *BucketManager) GetObject(ctx context.Context, item types.Keyable) error {
	err := Get(ctx, b.Bucket, item)
	return err
}

func (b *BucketManager) DeleteObject(ctx context.Context, item types.Keyable) error {
	err := Delete(ctx, b.Bucket, item)
	return err
}
