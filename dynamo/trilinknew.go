package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts all entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *TriLink[T0, T1, T2]) CheckLink(ctx context.Context, linkWrapper types.Linkable, entity0 T0, entity1 T1, entity2 T2) (linkExists bool, err error) {
	*link = *NewTriLink(entity0, entity1, entity2)
	// load the entities
	loaded0, err := link.Get(ctx, link.Entity0)
	if err != nil {
		return false, err
	}
	if !loaded0 {
		return false, &ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	loaded1, err := link.Get(ctx, link.Entity1)
	if err != nil {
		return false, err
	}
	if !loaded1 {
		return false, &ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	loaded2, err := link.Get(ctx, link.Entity2)
	if err != nil {
		return false, err
	}
	if !loaded2 {
		return false, &ErrEntityNotFound[T2]{Entity: link.Entity2}
	}

	linkExists, err = link.Get(ctx, linkWrapper)
	return linkExists, nil
}

// NewTriLink creates a new TriLink instance.
func NewTriLink[T0, T1, T2 types.Linkable](entity0 T0, entity1 T1, entity2 T2) *TriLink[T0, T1, T2] {
	link := TriLink[T0, T1, T2]{DiLink: DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}, Entity2: entity2}
	link.GenerateTriLinkCompositeKey()
	return &link
}
