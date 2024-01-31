package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"
)

func Delete(ctx context.Context, bucket string, item types.Keyable) error {
	client := clients.GetDefaultClient(ctx)
	pk, sk, err := item.Keys(0)
	if err != nil {
		return err
	}
	s3Key := pk + "/" + sk
	_, err = DeleteObjectWithClient(ctx, client, bucket, s3Key)
	return err
}

func DeleteObjectWithClient(ctx context.Context, client *clients.Client, bucket string, key string) (*s3.DeleteObjectOutput, error) {
	return client.S3().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}
