package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *DiLink[T0, T1]) CheckLink(ctx context.Context, linkWrapper DiLinkReturner[T0, T1]) (loaded bool, err error) {
	l := linkWrapper.ReturnLink()
	if l == nil {
		return false, nil
	}
	link = l
	return checkDiLink[T0, T1](ctx, *l)
}

// NewDiLink creates a new DiLink instance.
func NewDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) DiLink[T0, T1] {
	link := DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}
	link.GenerateDiLinkCompositeKey()
	return link
}

type DiLinkReturner[T0, T1 types.Linkable] interface {
	ReturnLink() (link *DiLink[T0, T1])
}

// CheckDiLink creates a new DiLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func checkDiLink[T0, T1 types.Linkable](ctx context.Context, link DiLink[T0, T1]) (allEntitiesExist bool, err error) {
	linkLoaded, err := link.Get(ctx, &link)
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
	if !loaded0 {
		return false, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if !loaded1 {
		return false, ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	if !linkLoaded {
		return true, ErrLinkNotFound{}
	}
	return true, nil
}
