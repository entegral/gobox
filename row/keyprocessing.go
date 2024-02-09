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
	PK       string
	SK       string
	Index    int
	IsEntity bool
}

type ErrInvalidKeySegment struct {
	label string
	value string
}

func (e ErrInvalidKeySegment) Error() string {
	return fmt.Sprintf("invalid key segment: %s(%s)", e.label, e.value)
}

// Define a function to generate the keys
func (item *Row[T]) GenerateKeys(ctx context.Context) (keys chan Key, errs chan error) {
	keys = make(chan Key, item.MaxGSIs())
	errs = make(chan error)
	defer close(keys)
	defer close(errs)

	for i := 0; i < item.MaxGSIs(); i++ {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return keys, errs
		default:
			key, err := item.object.Keys(i)
			if err != nil {
				errs <- fmt.Errorf("error generating keys for gsi %d of type %s: %w", i, item.object.Type(), err)
				return keys, errs
			}

			if key.PK == "" && key.SK == "" && i == 0 {
				// If the primary keys are empty and this is the primary index, assign a guid to the primary key and "default" to the sort key
				key.PK = uuid.UUIDv4()
				key.SK = "default"
			} else if key.PK == "" && key.SK == "" {
				continue
			} else if key.PK == "" {
				errs <- fmt.Errorf("partition key is required for gsi %d of type %s", i, item.object.Type())
				return keys, errs
			} else if key.SK == "" {
				errs <- fmt.Errorf("sort key is required for gsi %d of type %s", i, item.object.Type())
				return keys, errs
			}
			key, err = item.postProcessKey(ctx, key)
			if err != nil {
				errs <- fmt.Errorf("error post-processing keys for gsi %d of type %s: %w", i, item.object.Type(), err)
				return keys, errs
			}
			item.Keys.SetKey(key)
			keys <- key
		}
	}
	return keys, errs
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
	item.Keys.SetKey(Key{PK: pk, SK: sk, Index: key.Index, IsEntity: key.IsEntity})

	// Return the post-processed key
	return key, nil
}

func (item *Row[T]) postProcessKey(ctx context.Context, key Key) (processedKey Key, err error) {
	// If a custom KeysPostProcessor function is set, call it and return
	if item.KeysPostProcessor != nil {
		return item.KeysPostProcessor(ctx, key)
	} else {
		return item.DefaultPostProcessing(ctx, key)
	}
}

// Define a method to post-process the keys
func (item *Row[T]) postProcessKeys(ctx context.Context, keys chan Key) (<-chan Key, <-chan error) {
	processedKeys := make(chan Key)
	errs := make(chan error)
	// Close the processedKeys channel
	defer close(processedKeys)
	defer close(errs)

	// Otherwise, use the default post-processing logic
	for key := range keys {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return processedKeys, errs
		default:
			postProcessedKey, err := item.postProcessKey(ctx, key)
			if err != nil {
				errs <- err
				return processedKeys, errs
			}
			// Send the post-processed key to the processedKeys channel
			processedKeys <- postProcessedKey
		}
	}
	return processedKeys, errs
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
