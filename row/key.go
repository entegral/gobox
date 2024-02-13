package row

import (
	"context"
)

// Define a struct to hold the keys and the index
type Key struct {
	Pk       string `json:"pk" dynamodbav:"pk"`
	Sk       string `json:"sk" dynamodbav:"sk"`
	Index    int    `json:"index" dynamodbav:"index"`
	IsEntity bool   `json:"isEntity" dynamodbav:"-"`
}

// IndexName is a function to get the index name
func (key Key) IndexName() *string {
	i := GetIndexName(key)
	return &i
}

// Define a function to generate the keys
func (item *Row[T]) GenerateKeys(ctx context.Context) error {
	// Generate the keys for the row
	return item.Keys.GenerateKeys(item.object)
}
