package dynamo

import (
	"context"
	"sync"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/entegral/gobox/types"
)

type Result struct {
	Index  int
	Loaded bool
	Error  error
}

// BatchGet gets multiple rows from DynamoDB concurrently. The rows must implement the Keyable interface.
func (d *DBManager) BatchGet(ctx context.Context, rows []types.Linkable) <-chan Result {
	results := make(chan Result)

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row types.Linkable) {
				defer wg.Done()
				loaded, err := d.Get(ctx, row)
				results <- Result{i, loaded, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchPut puts multiple rows into DynamoDB concurrently. The rows must implement the Linkable interface.
func (d *DBManager) BatchPut(ctx context.Context, rows []types.Linkable) <-chan Result {
	results := make(chan Result)

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row types.Linkable) {
				defer wg.Done()
				err := d.Put(ctx, row)
				results <- Result{i, err == nil, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchDelete deletes multiple rows from DynamoDB concurrently. The rows must implement the Keyable interface.
func (d *DBManager) BatchDelete(ctx context.Context, rows []types.Linkable) <-chan Result {
	results := make(chan Result)

	go func() {
		var wg sync.WaitGroup

		for i, row := range rows {
			wg.Add(1)
			go func(i int, row types.Linkable) {
				defer wg.Done()
				err := d.Delete(ctx, row)
				results <- Result{i, err == nil, err}
			}(i, row)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

// BatchLoadFromMessage unmarshals multiple SQS messages into Rows and then loads the full items from DynamoDB concurrently.
func (d *DBManager) BatchLoadFromMessage(ctx context.Context, messages []sqstypes.Message, rows []types.Linkable) <-chan Result {
	results := make(chan Result)

	go func() {
		var wg sync.WaitGroup

		for i, message := range messages {
			wg.Add(1)
			go func(i int, message sqstypes.Message, row types.Linkable) {
				defer wg.Done()
				loaded, err := d.LoadFromMessage(ctx, message, row)
				results <- Result{i, loaded, err}
			}(i, message, rows[i])
		}

		wg.Wait()
		close(results)
	}()

	return results
}
