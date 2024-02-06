package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (item Row[T]) Get(ctx context.Context) (err error) {
	// Generate the keys from the input
	partitionKey, sortKey, err := item.object.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return err
	}

	// Create the key for DynamoDB
	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{
			Value: partitionKey,
		},
		"sk": &awstypes.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Create the GetItem input
	getItemInput := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(item.TableName()),
	}

	// Call DynamoDB GetItem
	result, err := item.GetClient(ctx).Dynamo().GetItem(ctx, getItemInput)
	if err != nil {
		return err
	}

	// Unmarshal the result into a Row
	return item.unmarshalMap(result.Item)
}
