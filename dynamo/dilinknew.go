package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *DiLink[T0, T1]) CheckLink(ctx context.Context, linkWrapper types.Linkable, entity0 T0, entity1 T1) (linkExists bool, err error) {
	*link = *NewDiLink(entity0, entity1)
	// load the entities
	loaded0, err := link.Get(ctx, link.Entity0)
	if err != nil {
		return false, err
	}
	if !loaded0 {
		return false, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	loaded1, err := link.Get(ctx, link.Entity1)
	if err != nil {
		return false, err
	}
	if !loaded1 {
		return false, ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	linkExists, err = link.Get(ctx, linkWrapper)
	return linkExists, err
}

// NewDiLink creates a new DiLink instance.
func NewDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) *DiLink[T0, T1] {
	link := DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}
	link.GenerateDiLinkKeys()
	return &link
}
