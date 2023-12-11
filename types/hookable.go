package types

import (
	"context"

	"gobox/clients"
)

// Hookable defines an object that has before and after hooks designed for use in
// situations where the item will be persisted
type Hookable interface {
	Before(ctx context.Context, clients *clients.Client) error
	After(ctx context.Context, clients *clients.Client) error
}
