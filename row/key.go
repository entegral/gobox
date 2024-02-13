package row

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Define a struct to hold the keys and the index
type Key struct {
	Pk               string                             `json:"pk"`
	Sk               string                             `json:"sk"`
	Index            int                                `json:"index"`
	IsEntity         bool                               `json:"isEntity"`
	LastEvaluatedKey map[string]awstypes.AttributeValue `json:"lastEvaluatedKey"`
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
