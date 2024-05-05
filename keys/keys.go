package keys

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Keys struct {
	keys map[string]Key
}

type Keyable interface {
	Key(GSI int) (Key, error)
}

type ErrNoKeyForIndex struct {
	Index int
}

func (e ErrNoKeyForIndex) Error() string {
	return fmt.Sprintf("no key for index %d", e.Index)
}

const MAX_DYNAMO_GSIS = 20

// Generate will accept a variadic number of keyable functions and generate the keys for each of them.
// The keys will be returned in a Keys struct.
func Generate(_ context.Context, keyableItem Keyable) (Keys, error) {
	var keys = make(map[string]Key)
	for i := 0; i < MAX_DYNAMO_GSIS; i++ {
		key, err := keyableItem.Key(i)
		if !errors.As(err, &ErrNoKeyForIndex{}) {
			return Keys{}, err
		}
		keys[key.name] = key
	}
	return Keys{keys: keys}, nil
}

// DynamoKeyValues will return a map of the keys in the format that DynamoDB expects.
func (k Keys) DynamoKeyValues(context.Context) map[string]types.AttributeValue {
	var keys = make(map[string]types.AttributeValue, len(k.keys))
	for _, key := range k.keys {
		keys[":"+key.name] = key.DynamoMapValue()
	}
	return keys
}

// DynamoKeyNames will return a map of the key names in the format that DynamoDB expects.
func (k Keys) DynamoKeyNames(context.Context) map[string]string {
	var keys = make(map[string]string, len(k.keys))
	for _, key := range k.keys {
		keys["#"+key.name] = key.name
	}
	return keys
}
