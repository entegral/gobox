package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		TableName: aws.String(item.TableName()),
	}

	// Call DynamoDB GetItem
	result, err := item.GetClient(ctx).Dynamo().GetItem(ctx, getItemInput)
	if err != nil {
		return err
	}

	// Unmarshal the result into a Row
	err = item.unmarshalMap(result.Item)
	if err != nil {
		return err
	}
	return nil
}

func (item *Row[T]) getBatchItems(ctx context.Context, keys []Key, results chan<- Row[T], errors chan<- error) {
	// Split the keys into batches of 100
	// batches := splitIntoBatches(keys, 100)

	// // Create a wait group to wait for all goroutines to finish
	// var wg sync.WaitGroup
	// wg.Add(len(batches))

	// // Process each batch concurrently
	// for _, batch := range batches {
	// 	go func(batch []Key) {
	// 		defer wg.Done()

	// 		// Create the BatchGetItem input
	// 		batchGetItemInput := item.createBatchGetItemInput(batch)

	// 		// Call DynamoDB BatchGetItem
	// 		result, err := item.GetClient(ctx).Dynamo().BatchGetItem(ctx, batchGetItemInput)
	// 		if err != nil {
	// 			errors <- err
	// 			return
	// 		}

	// 		// Unmarshal the results into Rows and send them to the results channel
	// 		for _, items := range result.Responses {
	// 			for _, itemMap := range items {
	// 				var row Row[T]
	// 				err = attributevalue.UnmarshalMap(itemMap, &row)
	// 				if err != nil {
	// 					errors <- err
	// 					return
	// 				}

	// 				results <- row
	// 			}
	// 		}
	// 	}(batch)
	// }

	// // Wait for all goroutines to finish
	// wg.Wait()
}

// // splitIntoBatches splits a slice of keys into batches of the specified size.
// func splitIntoBatches(keys []Key, batchSize int) [][]Key {
// 	var batches [][]Key

// 	for batchSize < len(keys) {
// 		keys, batches = keys[batchSize:], append(batches, keys[0:batchSize:batchSize])
// 	}

// 	batches = append(batches, keys)

// 	return batches
// }

func (item *Row[T]) createBatchGetItemInput(keys []Key) *dynamodb.BatchGetItemInput {
	keysMap := make(map[string]types.AttributeValue)

	for _, key := range keys {
		keysMap = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: key.PK,
			},
			"sk": &types.AttributeValueMemberS{
				Value: key.SK,
			},
		}
	}

	keysAndAttributes := types.KeysAndAttributes{
		Keys: []map[string]types.AttributeValue{keysMap},
	}

	// Create the BatchGetItem input
	return &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			item.TableName(): keysAndAttributes,
		},
	}
}
