package types

import (
	"context"

	"github.com/entegral/gobox/clients"
)

// Hookable defines an object that has before and after hooks designed for use in
// situations where the item will be persisted
type Hookable interface {
	Before(ctx context.Context, clients *clients.Client) error
	After(ctx context.Context, clients *clients.Client) error
}
