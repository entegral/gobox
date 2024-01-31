package s3

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"
)

func Get(ctx context.Context, bucket string, item types.Keyable) error {
	client := clients.GetDefaultClient(ctx)
	pk, sk, err := item.Keys(0)
	if err != nil {
		return err
	}
	s3Key := pk + "/" + sk
	_, err = GetObjectWithClient(ctx, client, bucket, s3Key, item)
	return err
}

func GetObjectWithClient(ctx context.Context, client *clients.Client, bucket string, key string, item types.Keyable) (*s3.GetObjectOutput, error) {
	out, err := client.S3().GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	body, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, item)
	if err != nil {
		return nil, err
	}
	return out, nil
}
