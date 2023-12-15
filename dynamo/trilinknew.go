package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts all entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *TriLink[T0, T1, T2]) CheckLink(ctx context.Context, linkWrapper types.Linkable, entity0 T0, entity1 T1, entity2 T2) (allEntitiesExist bool, err error) {
	*link = NewTriLink[T0, T1, T2](link.Entity0, link.Entity1, link.Entity2)
	return checkTriLink[T0, T1, T2](
		ctx,
		linkWrapper,
		link.Entity0,
		link.Entity1,
		link.Entity2,
	)
}

// NewTriLink creates a new TriLink instance.
func NewTriLink[T0, T1, T2 types.Linkable](entity0 T0, entity1 T1, entity2 T2) TriLink[T0, T1, T2] {
	link := TriLink[T0, T1, T2]{DiLink: DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}, Entity2: entity2}
	link.GenerateTriLinkCompositeKey()
	return link
}

// CheckTriLink creates a new TriLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func checkTriLink[T0, T1, T2 types.Linkable](ctx context.Context, triLinkWrapper types.Linkable, entity0 T0, entity1 T1, entity2 T2) (allEntitiesExist bool, err error) {
	link := NewTriLink[T0, T1, T2](entity0, entity1, entity2)
	linkLoaded, err := link.Get(ctx, triLinkWrapper)
	if err != nil {
		return false, err
	}
	// load the entities
	loaded0, err := link.Get(ctx, link.Entity0)
	if err != nil {
		return false, err
	}
	loaded1, err := link.Get(ctx, link.Entity1)
	if err != nil {
		return false, err
	}
	loaded2, err := link.Get(ctx, link.Entity2)
	if err != nil {
		return false, err
	}
	if !loaded0 {
		return false, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if !loaded1 {
		return false, ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	if !loaded2 {
		return false, ErrEntityNotFound[T2]{Entity: link.Entity2}
	}
	if !linkLoaded {
		return true, ErrLinkNotFound{}
	}
	return true, nil
}
