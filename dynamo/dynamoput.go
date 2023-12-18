package dynamo

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// PutItem puts a row into DynamoDB. The row must implement the
// Keyable interface. This method uses the default client. If you need to use a specific
// client, use PutItemWithClient instead, or use the client.SetDefaultClient method.
func PutItem(ctx context.Context, row types.Linkable) (*dynamodb.PutItemOutput, error) {
	tn := clients.TableName(ctx)
	return putItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), tn, row)
}

func PutItemWithTablename(ctx context.Context, tablename string, row types.Linkable) (*dynamodb.PutItemOutput, error) {
	return putItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), tablename, row)
}

func putItemPrependTypeWithClient(ctx context.Context, client *clients.Client, tablename string, row types.Linkable) (*dynamodb.PutItemOutput, error) {
	pk, sk, err := row.Keys(0)
	if err != nil {
		return nil, err
	}
	av, err := attributevalue.MarshalMap(row)
	if err != nil {
		return nil, err
	}
	pkWithTypePrefix, err := addKeySegment(rowType, row.Type())
	if err != nil {
		return nil, err
	}
	seg, err := addKeySegment(rowPk, pk)
	if err != nil {
		return nil, err
	}
	pkWithTypePrefix += seg
	av["pk"] = &awstypes.AttributeValueMemberS{Value: pkWithTypePrefix}
	av["sk"] = &awstypes.AttributeValueMemberS{Value: sk}
	av["type"] = &awstypes.AttributeValueMemberS{Value: row.Type()}
	return putItemWithClient(ctx, client, tablename, av)
}

// putItemWithClient puts a row into DynamoDB using the provided client.
func putItemWithClient(ctx context.Context, client *clients.Client, tablename string, av map[string]awstypes.AttributeValue) (*dynamodb.PutItemOutput, error) {
	return client.Dynamo().PutItem(ctx, &dynamodb.PutItemInput{
		TableName:    &tablename,
		Item:         av,
		ReturnValues: awstypes.ReturnValueAllOld,
	})
}

// Shardable is an interface that can be implemented by a row to indicate that it
// should be sharded when saved to dynamo. The Shard method should return a string
// that will be appended to the pk to create the final pk.
type Shardable interface {
	types.Keyable
	MaxShard() int
}

func getShard(maxShard int) string {
	shard := rand.Intn(maxShard)
	return fmt.Sprintf(".%d", shard)
}

// PutItemWithShard puts a row into DynamoDB using the provided client and shard.
func PutItemWithShard(ctx context.Context, client *clients.Client, row Shardable) (*dynamodb.PutItemOutput, error) {
	pk, sk, err := row.Keys(0)
	if err != nil {
		return nil, err
	}
	av, err := attributevalue.MarshalMap(row)
	if err != nil {
		return nil, err
	}
	pk = pk + getShard(row.MaxShard())
	av["pk"] = &awstypes.AttributeValueMemberS{Value: pk}
	av["sk"] = &awstypes.AttributeValueMemberS{Value: sk}
	return putItemWithClient(ctx, client, client.TableName(ctx), av)
}
