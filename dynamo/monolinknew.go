package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entity itself, only the link.
func (link *MonoLink[T0]) CheckLink(ctx context.Context, linkWrapper types.Typeable, input T0) (linkLoaded bool, err error) {
	*link = NewMonoLink[T0](input)
	link.UnmarshalledType = linkWrapper.Type()
	linkLoaded, err = link.Get(ctx, link)
	if err != nil {
		return false, err
	}
	if !linkLoaded {
		return false, ErrLinkNotFound{}
	}
	return true, nil
}

// NewMonoLink creates a new MonoLink instance.
func NewMonoLink[T0 types.Linkable](entity0 T0) MonoLink[T0] {
	link := MonoLink[T0]{Entity0: entity0}
	link.GenerateMonoLinkCompositeKey()
	return link
}

// CheckMonoLink creates a new MonoLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func CheckMonoLink[T0 types.Linkable](entity0 T0) (*MonoLink[T0], error) {
	link := NewMonoLink[T0](entity0)
	linkLoaded, err := link.Get(context.Background(), &link)
	if err != nil {
		return &link, err
	}
	// load the entities
	loaded0, err := link.Get(context.Background(), link.Entity0)
	if err != nil {
		return &link, err
	}
	if !loaded0 {
		return &link, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if !linkLoaded {
		return &link, ErrLinkNotFound{}
	}
	return &link, nil
}
