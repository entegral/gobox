package s3

import (
	"context"

	"github.com/entegral/gobox/types"
)

type BucketManager struct {
	Bucket string
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
