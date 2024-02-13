package row

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ItemSet[T Rowable] struct {
	Items []Row[T]
}

func (is *ItemSet[T]) GenerateKeys(ctx context.Context) (<-chan Key, <-chan error) {
	keys := make(chan Key)
	errs := make(chan error)

	go func() {
		defer close(keys)
		defer close(errs)

		for _, item := range is.Items {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				keygenErrs := item.GenerateKeys(ctx)
				if keygenErrs != nil {
					errs <- keygenErrs
					return
				}
			}
		}
	}()

	return keys, errs
}

type PutResult[T Rowable] struct {
	UnprocessedItems []Row[T]
}

func (set *ItemSet[T]) BatchPut(ctx context.Context) (<-chan PutResult[T], <-chan error) {
	// Create channels for results and errors
	results := make(chan PutResult[T])
	errs := make(chan error)

	go func() {
		defer close(results)
		defer close(errs)

		// Split the items into batches of 25
		batches := splitIntoBatches(set.Items, 25)

		// Create a wait group to wait for all goroutines to finish
		var wg sync.WaitGroup
		wg.Add(len(batches))

		// Process each batch concurrently
		for _, batch := range batches {
			go func(batch []Row[T]) {
				defer wg.Done()

				// Prepare the BatchWriteItemInput
				writeRequests := make([]types.WriteRequest, len(batch))
				for i, item := range batch {
					itemData, err := attributevalue.MarshalMap(item)
					if err != nil {
						errs <- err
						return
					}
					writeRequests[i] = types.WriteRequest{
						PutRequest: &types.PutRequest{
							Item: itemData,
						},
					}
				}
				tablename := batch[0].TableName()
				if tablename == nil {
					errs <- errors.New("TableName is nil for item: " + batch[0].Type())
					return
				}
				input := &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						*tablename: writeRequests,
					},
				}

				// Call DynamoDB BatchWriteItem
				resp, err := batch[0].GetClient(ctx).Dynamo().BatchWriteItem(ctx, input)
				if err != nil {
					errs <- err
					return
				}

				// Handle the response
				var unprocessedItems []Row[T]
				for _, writeRequests := range resp.UnprocessedItems {
					for _, writeRequest := range writeRequests {
						if writeRequest.PutRequest != nil {
							var item Row[T]
							err := attributevalue.UnmarshalMap(writeRequest.PutRequest.Item, &item)
							if err != nil {
								errs <- err
								return
							}
							unprocessedItems = append(unprocessedItems, item)
						}
					}
				}

				results <- PutResult[T]{UnprocessedItems: unprocessedItems}
			}(batch)
		}

		// Wait for all goroutines to finish
		wg.Wait()
	}()

	return results, errs
}

// splitIntoBatches splits a slice of Row[T] into batches of the specified size.
func splitIntoBatches[T Rowable](items []Row[T], batchSize int) [][]Row[T] {
	var batches [][]Row[T]

	for batchSize < len(items) {
		items, batches = items[batchSize:], append(batches, items[0:batchSize:batchSize])
	}

	if len(items) > 0 {
		batches = append(batches, items)
	}

	return batches
}
