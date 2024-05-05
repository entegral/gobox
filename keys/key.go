package keys

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Key struct {
	name  string
	value string
	index int
}

// NewKey creates a new key with the given name and value and index.
func NewKey(name, value string, index int) Key {
	return Key{name, value, index}
}

func (k Key) DynamoMapValue() types.AttributeValue {
	return &types.AttributeValueMemberS{
		Value: k.value,
	}
}

func (k Key) DynamoMapKey() string {
	return "#" + k.name
}
