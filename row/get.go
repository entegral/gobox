package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Generates keys using the object and attempts to get it from DynamoDB
func (item *Row[T]) Get(ctx context.Context) (err error) {
	// Generate the key from the input
	partitionKey, sortKey, err := item.object.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return err
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
		return <-errs
	}

	// Create the key for DynamoDB
	dynamoKey := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{
			Value: processedKey.PK,
		},
		"sk": &awstypes.AttributeValueMemberS{
			Value: processedKey.SK,
		},
	}

	// Create the GetItem input
	getItemInput := &dynamodb.GetItemInput{
		Key:       dynamoKey,
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
