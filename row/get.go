package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Get retrieves rows from DynamoDB that match the provided keys.
// If a single key is provided, it retrieves the row and unmarshals it into the receiver.
// If multiple keys are provided, it retrieves the rows in batches and returns them through a channel.
func (item *Row[T]) Get(ctx context.Context) error {
	return item.getSingleItem(ctx)
}

func (item *Row[T]) getSingleItem(ctx context.Context) error {
	item.GenerateKeys(ctx)
	// Create the GetItem input
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: item.Keys.Pk,
			},
			"sk": &types.AttributeValueMemberS{
				Value: item.Keys.Sk,
			},
		},
		TableName: item.TableName(),
	}

	// Call DynamoDB GetItem
	result, err := item.GetClient(ctx).Dynamo().GetItem(ctx, getItemInput)
	if err != nil {
		return err
	}

	// Unmarshal the result into a Row
	return item.unmarshalMap(result.Item)
}
