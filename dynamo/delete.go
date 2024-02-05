package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Delete[T Rowable](ctx context.Context, item T) (oldRow row[T], err error) {
	r := newRow[T](item)
	// Generate the keys from the input
	partitionKey, sortKey, err := r.object.Keys(0) // Assuming GSI 0 for primary key
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
		TableName:    aws.String(r.TableName()),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB DeleteItem
	result, err := r.GetClient(ctx).Dynamo().DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}
