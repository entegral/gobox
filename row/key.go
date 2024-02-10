package row

import (
	"context"
)

// Define a function to generate the keys
func (item *Row[T]) GenerateKeys(ctx context.Context) error {
	// Generate the keys for the row
	return item.Keys.GenerateKeys(item.object)
}
