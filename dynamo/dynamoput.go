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
	return putItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), row)
}

func putItemPrependTypeWithClient(ctx context.Context, client *clients.Client, row types.Linkable) (*dynamodb.PutItemOutput, error) {
	pk, sk, err := row.Keys(0)
	if err != nil {
		return nil, err
	}
	av, err := attributevalue.MarshalMap(row)
	if err != nil {
		return nil, err
	}
	pkWithTypePrefix, err := prependWithRowType(row, pk)
	if err != nil {
		return nil, err
	}
	av["pk"] = &awstypes.AttributeValueMemberS{Value: pkWithTypePrefix}
	av["sk"] = &awstypes.AttributeValueMemberS{Value: sk}
	av["type"] = &awstypes.AttributeValueMemberS{Value: row.Type()}
	tn := row.TableName(ctx)
	return putItemWithClient(ctx, client, tn, av)
}

// putItemWithClient puts a row into DynamoDB using the provided client.
func putItemWithClient(ctx context.Context, client *clients.Client, tablename string, av map[string]awstypes.AttributeValue) (*dynamodb.PutItemOutput, error) {
	rcc := awstypes.ReturnConsumedCapacityNone
	if checkTesting() {
		rcc = awstypes.ReturnConsumedCapacityTotal
	}
	return client.Dynamo().PutItem(ctx, &dynamodb.PutItemInput{
		TableName:              &tablename,
		Item:                   av,
		ReturnValues:           awstypes.ReturnValueAllOld,
		ReturnConsumedCapacity: rcc,
	})
}

// Shardable is an interface that can be implemented by a row to indicate that it
// should be sharded when saved to dynamo. The Shard method should return a string
// that will be appended to the pk to create the final pk.
type Shardable interface {
	types.Linkable
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
	pkWithShard := pk + getShard(row.MaxShard())
	prefixedPkWithShard, err := prependWithRowType(row, pkWithShard)
	av["pk"] = &awstypes.AttributeValueMemberS{Value: prefixedPkWithShard}
	av["sk"] = &awstypes.AttributeValueMemberS{Value: sk}
	tn := row.TableName(ctx)
	return putItemWithClient(ctx, client, tn, av)
}
