package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// GetFrom attempts to load the link associated with the given entity.
func (link *MonoLink[T0]) GetFrom(ctx context.Context, linkWrapper types.Typeable, input T0) (linkLoaded bool, err error) {
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

// CheckDiLink creates a new DiLink instance and attempts to load the entities.
// If either of the entities cannot be loaded from dynamo, an error will be returned.
//
// If you need a DiLink instance that does not require the entities to be loaded,
// you can use the NewDiLink function instead.
//
// If the link itself does not exist, an ErrLinkNotFound error will be returned,
// but the entities will still be loaded and you can call the .Link() method to
// create the link in dynamo.
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
