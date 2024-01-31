package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"
)

func Put(ctx context.Context, bucket string, item types.Keyable) error {
	client := clients.GetDefaultClient(ctx)
	pk, sk, err := item.Keys(0)
	if err != nil {
		return err
	}
	s3Key := pk + "/" + sk
	_, err = PutObjectWithClient(ctx, client, bucket, s3Key, item)
	return err
}

func PutObjectWithClient(ctx context.Context, client *clients.Client, bucket string, key string, item any) (*s3.PutObjectOutput, error) {
	reader, err := ToReader(item)
	if err != nil {
		return nil, err
	}
	return client.S3().PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   reader,
	})
}

func ToReader(item any) (io.Reader, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
