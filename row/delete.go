package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeleteResult[T Rowable] struct {
	DeletedItems    []Row[T]
	UnprocessedKeys []Key
}

// Delete removes items from DynamoDB identified by the provided keys.
// If a single key is provided, it deletes the item directly.
// If multiple keys are provided, it deletes the items in batches.
// It returns two channels: one for the delete results and one for errors.
// Each DeleteResult contains the deleted items and any unprocessed keys.
// Errors from individual delete operations are sent to the error channel.
//
// Note: The caller is responsible for handling and draining both channels.
func (item *Row[T]) Delete(ctx context.Context) error {
	return item.deleteSingleItem(ctx)
}

func (item *Row[T]) deleteSingleItem(ctx context.Context) error {
	item.GenerateKeys(ctx)

	// Create the DeleteItem input
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: item.Keys.Pk,
			},
			"sk": &types.AttributeValueMemberS{
				Value: item.Keys.Sk,
			},
		},
		TableName:    aws.String(item.TableName()),
		ReturnValues: types.ReturnValueAllOld, // Return all attributes of the deleted item
	}

	// Call DynamoDB DeleteItem
	result, err := item.GetClient(ctx).Dynamo().DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return err
	}

	// Unmarshal the result into a Row
	return item.unmarshalMap(result.Attributes)
}
