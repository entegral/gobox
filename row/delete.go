package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (item Row[T]) Delete(ctx context.Context) (oldRow Row[T], err error) {
	// Generate the keys from the input
	partitionKey, sortKey, err := item.object.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return oldRow, err
	}

	// Create the key for DynamoDB
	key := map[string]awstypes.AttributeValue{
		"pk0": &awstypes.AttributeValueMemberS{
			Value: partitionKey,
		},
		"sk0": &awstypes.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Create the DeleteItem input
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key:          key,
		TableName:    aws.String(item.TableName()),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB DeleteItem
	result, err := item.GetClient(ctx).Dynamo().DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}
