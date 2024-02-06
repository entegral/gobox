package row

import (
	"context"
	"fmt"
)

// Define a struct to hold the keys and the index
type Key struct {
	PK    string
	SK    string
	Index int
}

// Define a function to generate the keys
func (item Row[T]) GenerateKeys(ctx context.Context, keys chan<- Key, errs chan<- error) {
	defer close(keys)
	defer close(errs)

	for i := 0; i < item.MaxGSIs(); i++ {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return
		default:
			pk, sk, err := item.object.Keys(i)
			if err != nil {
				errs <- fmt.Errorf("error generating keys for gsi %d of type %s: %w", i, item.object.Type(), err)
				return
			}

			if pk == "" && sk == "" {
				continue
			} else if pk == "" {
				errs <- fmt.Errorf("partition key is required for gsi %d of type %s", i, item.object.Type())
				return
			} else if sk == "" {
				errs <- fmt.Errorf("sort key is required for gsi %d of type %s", i, item.object.Type())
				return
			}

			keys <- Key{PK: pk, SK: sk, Index: i}
		}
	}
}

// Define the default post-processing function
func (item *Row[T]) DefaultPostProcessing(ctx context.Context, key Key) (Key, error) {
	// Default post-processing logic goes here
	// ...

	// If an error occurs during post-processing, return it
	// return Key{}, err

	// Return the post-processed key
	return key, nil
}

// Define a method to post-process the keys
func (item *Row[T]) PostProcessKeys(ctx context.Context, keys <-chan Key, processedKeys chan<- Key, errs chan<- error) {
	// Close the processedKeys channel
	defer close(processedKeys)

	// If a custom KeysPostProcessor function is set, call it and return
	if item.KeysPostProcessor != nil {
		item.KeysPostProcessor(ctx, keys, processedKeys, errs)
		return
	}

	// Otherwise, use the default post-processing logic
	for key := range keys {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return
		default:
			// Call the default post-processing function
			postProcessedKey, err := item.DefaultPostProcessing(ctx, key)
			if err != nil {
				errs <- err
				return
			}

			// Send the post-processed key to the processedKeys channel
			processedKeys <- postProcessedKey
		}
	}
}
