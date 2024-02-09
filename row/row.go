package row

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/types"
)

type Rowable interface {
	Keyable
	types.Typeable
}

// KeyPostProcessor is a function that allows for post-processing of keys
type KeyPostProcessor func(ctx context.Context, key Key) (processedKey Key, err error)

type Row[T Rowable] struct {
	// Table configuration, marshaled to json but not to dynamo
	Table `dynamodbav:"-" json:"tablename,omitempty"`

	// The keys for the row
	Keys

	// KeysPostProcessor is a function that allows for post-processing of keys
	// if this is set, it will be called instead of the default PostProcessor
	KeysPostProcessor KeyPostProcessor

	// The object that is being stored or retrieved
	object T
}

func NewRow[T Rowable](object T) Row[T] {
	return Row[T]{object: object}
}

func (r *Row[T]) Type() string {
	return r.object.Type()
}

func (r *Row[T]) unmarshalMap(m map[string]awstypes.AttributeValue) error {
	// Create a new map to hold the non-key values
	err := r.Keys.unmarshalKeysFromMap(m)
	if err != nil {
		return err
	}

	// Unmarshal the non-key values into the Object
	err = attributevalue.UnmarshalMap(m, &r.object)
	return err
}

// Object returns the object that is being stored or retrieved
func (r *Row[T]) Object() T {
	return r.object
}

// UnmarshalJSON unmarshals the JSON data into the Row.
// It expects the JSON data to be a map with the following keys:
// - Table: the table configuration
// - Keys: the keys for the row
// - Object: the object that is being stored or retrieved
func (r *Row[T]) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into a map
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(3)

	// Use a channel to collect errors from the goroutines
	errs := make(chan error, 3)

	// Unmarshal the table configuration in a goroutine
	go func() {
		defer wg.Done()
		err := json.Unmarshal(m["Table"], &r.Table)
		if err != nil {
			errs <- err
		}
	}()

	// Unmarshal the keys in a goroutine
	go func() {
		defer wg.Done()
		err := json.Unmarshal(m["Keys"], &r.Keys)
		if err != nil {
			errs <- err
		}
	}()

	// Unmarshal the object in a goroutine
	go func() {
		defer wg.Done()
		err := json.Unmarshal(m["Object"], r.object)
		if err != nil {
			errs <- err
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the error channel and check if there were any errors
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Row[T]) MarshalJSON() ([]byte, error) {
	// Define a temporary struct that has the same structure as the Row type
	type tempRow struct {
		Table  Table           `json:"Table"`
		Keys   Keys            `json:"Keys"`
		Object json.RawMessage `json:"Object"`
	}

	// Marshal the object field into a raw JSON message
	objectJSON, err := json.Marshal(r.object)
	if err != nil {
		return nil, err
	}

	// Create a temporary struct with the fields of the Row type
	temp := tempRow{
		Table:  r.Table,
		Keys:   r.Keys,
		Object: objectJSON,
	}

	// Marshal the temporary struct into a JSON object
	return json.Marshal(temp)
}
