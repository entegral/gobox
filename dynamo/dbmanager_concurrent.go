package dynamo

import (
	"context"
	"sync"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/entegral/gobox/types"
)

type Result[T types.Linkable] struct {
	Index int
	Item  *T
	Error error
}

// BatchGet gets multiple rows from DynamoDB concurrently. The rows must implement the Keyable interface.
func (d DynamoManager[T]) BatchGet(ctx context.Context, rows []T) <-chan Result[T] {
	results := make(chan Result[T])

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row T) {
				defer wg.Done()
				item, err := d.Get(ctx)
				if err != nil {
					results <- Result[T]{i, item, err}
					return
				}
				results <- Result[T]{i, item, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchPut puts multiple rows into DynamoDB concurrently. The rows must implement the Linkable interface.
func (d *DynamoManager[T]) BatchPut(ctx context.Context, rows []T) <-chan Result[T] {
	results := make(chan Result[T])

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row T) {
				defer wg.Done()
				err := d.Put(ctx, row)

				results <- Result[T]{i, &row, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchDelete deletes multiple rows from DynamoDB concurrently. The rows must implement the Keyable interface.
func (d *DynamoManager[T]) BatchDelete(ctx context.Context, rows []T) <-chan Result[T] {
	results := make(chan Result[T])

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row T) {
				defer wg.Done()
				err := d.Delete(ctx)
				results <- Result[T]{i, &row, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchLoadFromMessage unmarshals multiple SQS messages into Rows and then loads the full items from DynamoDB concurrently.
func (d *DynamoManager[T]) BatchLoadFromMessage(ctx context.Context, messages []sqstypes.Message, rows []T) <-chan Result[T] {
	results := make(chan Result[T])

	go func() {
		var wg sync.WaitGroup

		for i, message := range messages {
			wg.Add(1)
			go func(i int, message sqstypes.Message, row T) {
				defer wg.Done()
				item, err := d.LoadFromMessage(ctx, message)
				results <- Result[T]{i, item, err}
			}(i, message, rows[i])
		}

		wg.Wait()
		close(results)
	}()

	return results
}
