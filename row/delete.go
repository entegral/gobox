package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (item *Row[T]) Delete(ctx context.Context) (oldRow Row[T], err error) {
	// Generate the key from the input
	partitionKey, sortKey, err := item.object.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return oldRow, err
	}

	// Create the Key
	key := Key{PK: partitionKey, SK: sortKey, Index: 0}

	// Create channels for keys and processed keys
	keys := make(chan Key, 1)
	processedKeys := make(chan Key, 1)
	errs := make(chan error, 1)

	// Send the key to the keys channel
	keys <- key
	close(keys)

	// Start a goroutine to post-process the key
	go item.postProcessKeys(ctx, keys, processedKeys, errs)

	// Receive the post-processed key from the processedKeys channel
	processedKey, ok := <-processedKeys
	if !ok {
		// If the processedKeys channel was closed without sending a key, receive the error from the errs channel
		return oldRow, <-errs
	}

	// ensure the item's key struct is updated with the new keys
	item.Keys.SetKey(processedKey)

	// Create the key for DynamoDB
	dynamoKey := map[string]awstypes.AttributeValue{
		"pk0": &awstypes.AttributeValueMemberS{
			Value: processedKey.PK,
		},
		"sk0": &awstypes.AttributeValueMemberS{
			Value: processedKey.SK,
		},
	}

	// Create the DeleteItem input
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key:          dynamoKey,
		TableName:    aws.String(item.TableName()),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB DeleteItem
	result, err := item.GetClient(ctx).Dynamo().DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	return oldRow, attributevalue.UnmarshalMap(result.Attributes, &oldRow)
}
