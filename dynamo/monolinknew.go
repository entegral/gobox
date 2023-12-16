package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entity itself, only the link.
func (link *MonoLink[T0]) CheckLink(ctx context.Context, linkWrapper MonoLinkReturner[T0]) (linkLoaded bool, err error) {
	l := linkWrapper.ReturnLink()
	if l == nil {
		return false, nil
	}
	link = l
	return checkMonoLink[T0](ctx, *l)
}

// NewMonoLink creates a new MonoLink instance.
func NewMonoLink[T0 types.Linkable](entity0 T0) MonoLink[T0] {
	link := MonoLink[T0]{Entity0: entity0}
	link.GenerateMonoLinkCompositeKey()
	return link
}

func (link *MonoLink[T0]) ReturnLink() *MonoLink[T0] {
	return link
}

type MonoLinkReturner[T0 types.Linkable] interface {
	ReturnLink() (link *MonoLink[T0])
}

// CheckMonoLink creates a new MonoLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func checkMonoLink[T0 types.Linkable](ctx context.Context, link MonoLink[T0]) (allEntitiesExist bool, err error) {
	linkLoaded, err := link.Get(ctx, &link)
	if err != nil {
		return false, err
	}
	// load the entities
	loaded0, err := link.Get(ctx, link.Entity0)
	if err != nil {
		return false, err
	}
	if !loaded0 {
		return false, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if !linkLoaded {
		return true, ErrLinkNotFound{}
	}
	return true, nil
}
