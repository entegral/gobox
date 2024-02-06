package row

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/dgryski/trifles/uuid"
	"github.com/entegral/gobox/types"
)

// Define a struct to hold the keys and the index
type Key struct {
	PK    string
	SK    string
	Index int
}

type ErrInvalidKeySegment struct {
	label string
	value string
}

func (e ErrInvalidKeySegment) Error() string {
	return fmt.Sprintf("invalid key segment: %s(%s)", e.label, e.value)
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

			if pk == "" && sk == "" && i == 0 {
				// If the primary keys are empty and this is the primary index, assign a guid to the primary key and "default" to the sort key
				pk = uuid.UUIDv4()
				sk = "default"
			} else if pk == "" && sk == "" {
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
	pk, err := prependWithRowType(item.object, key.PK)
	if err != nil {
		return key, err
	}

	sk, err := prependWithRowType(item.object, key.SK)
	if err != nil {
		return key, err
	}

	// Set the post-processed keys
	key.PK = pk
	key.SK = sk

	// Return the post-processed key
	return key, nil
}

// Define a method to post-process the keys
func (item *Row[T]) postProcessKeys(ctx context.Context, keys <-chan Key, processedKeys chan<- Key, errs chan<- error) {
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

func containsObscureWhitespace(value string) bool {
	for _, r := range value {
		if unicode.IsSpace(r) && !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

func addKeySegment(label linkLabels, value string) (string, error) {
	// Check if label or value contains characters that could affect the regex
	if len(value) == 0 || strings.ContainsAny(string(label), "()") || containsObscureWhitespace(value) {
		return "", ErrInvalidKeySegment{string(label), value}
	}
	if !label.IsValidLabel() {
		return "", ErrInvalidKeySegment{string(label), value}
	}

	// Check if value matches any linkLabel
	err := label.IsValidValue(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s(%s)", label, value), nil
}

func prependWithRowType(row types.Typeable, pk string) (string, error) {
	pkWithTypePrefix, err := addKeySegment(rowType, row.Type())
	if err != nil {
		return "", err
	}
	seg, err := addKeySegment(rowPk, pk)
	if err != nil {
		return "", err
	}
	pkWithTypePrefix += seg
	return pkWithTypePrefix, nil
}
