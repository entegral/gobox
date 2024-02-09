package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	key, err := item.object.Keys(0)
	if err != nil {
		return err
	}
	// Create the DeleteItem input
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: key.PK,
			},
			"sk": &types.AttributeValueMemberS{
				Value: key.SK,
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
	err = attributevalue.UnmarshalMap(result.Attributes, &item)
	if err != nil {
		return err
	}

	return nil
}

func (item *Row[T]) deleteBatchItems(ctx context.Context, keys []Key) (<-chan DeleteResult[T], <-chan error) {
	// Create channels for results and errors
	// results := make(chan DeleteResult[T])
	// errors := make(chan error)

	// go func() {
	// 	defer close(results)
	// 	defer close(errors)

	// 	// Split the keys into batches of 25
	// 	batches := splitIntoBatches(keys, 25)

	// 	// Create a wait group to wait for all goroutines to finish
	// 	var wg sync.WaitGroup
	// 	wg.Add(len(batches))

	// 	// Process each batch concurrently
	// 	for _, batch := range batches {
	// 		go func(batch []Key) {
	// 			defer wg.Done()

	// 			// Create the BatchWriteItem input
	// 			batchWriteItemInput := item.createBatchWriteItemInput(batch)

	// 			// Call DynamoDB BatchWriteItem
	// 			result, err := item.GetClient(ctx).Dynamo().BatchWriteItem(ctx, batchWriteItemInput)
	// 			if err != nil {
	// 				errors <- err
	// 				return
	// 			}

	// 			// If there are unprocessed items, send their keys to the results channel
	// 			unprocessedKeys := make([]Key, 0)
	// 			for _, writeRequest := range result.UnprocessedItems[item.TableName()] {
	// 				unprocessedKeys = append(unprocessedKeys, Key{
	// 					PK: writeRequest.DeleteRequest.Key["pk"].(*types.AttributeValueMemberS).Value,
	// 					SK: writeRequest.DeleteRequest.Key["sk"].(*types.AttributeValueMemberS).Value,
	// 				})
	// 			}
	// 			results <- DeleteResult[T]{UnprocessedKeys: unprocessedKeys}
	// 		}(batch)
	// 	}

	// 	// Wait for all goroutines to finish
	// 	wg.Wait()
	// }()

	// return results, errors
	return nil, nil
}

func (item *Row[T]) createBatchWriteItemInput(keys []Key) *dynamodb.BatchWriteItemInput {
	writeRequests := make([]types.WriteRequest, len(keys))

	for i, key := range keys {
		writeRequests[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{
						Value: key.PK,
					},
					"sk": &types.AttributeValueMemberS{
						Value: key.SK,
					},
				},
			},
		}
	}

	// Create the BatchWriteItem input
	return &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			item.TableName(): writeRequests,
		},
	}
}
