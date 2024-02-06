package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ItemSet[T Rowable] struct {
	Items []Row[T]
}

type PutResult[T Rowable] struct {
	UnprocessedItems []Row[T]
}

func (set *ItemSet[T]) Puts(ctx context.Context) (<-chan PutResult[T], <-chan error) {
	results := make(chan PutResult[T])
	errors := make(chan error)

	go func() {
		defer close(results)
		defer close(errors)

		// Prepare the BatchWriteItemInput
		writeRequests := make([]types.WriteRequest, len(set.Items))
		for i, item := range set.Items {
			itemData, err := attributevalue.MarshalMap(item)
			if err != nil {
				errors <- err
				return
			}
			writeRequests[i] = types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: itemData,
				},
			}
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				set.Items[0].TableName(): writeRequests,
			},
		}

		// Call DynamoDB BatchWriteItem
		resp, err := set.Items[0].GetClient(ctx).Dynamo().BatchWriteItem(ctx, input)
		if err != nil {
			errors <- err
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
						errors <- err
						return
					}
					unprocessedItems = append(unprocessedItems, item)
				}
			}
		}

		results <- PutResult[T]{UnprocessedItems: unprocessedItems}
	}()

	return results, errors
}
