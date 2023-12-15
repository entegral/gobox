package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *DiLink[T0, T1]) CheckLink(ctx context.Context, linkWrapper types.Typeable, entity0 T0, entity1 T1) (loaded bool, err error) {
	*link = NewDiLink[T0, T1](entity0, entity1)
	link.UnmarshalledType = linkWrapper.Type()
	loaded, err = link.Get(ctx, link)
	if err != nil {
		return false, err
	}
	if !loaded {
		return false, ErrLinkNotFound{}
	}
	return true, nil
}

// NewDiLink creates a new DiLink instance.
func NewDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) DiLink[T0, T1] {
	link := DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}
	link.GenerateDiLinkCompositeKey()
	return link
}

// CheckDiLink creates a new DiLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func CheckDiLink[T0, T1 types.Linkable](diLinkWrapper types.Linkable, entity0 T0, entity1 T1) (allEntitiesExist bool, err error) {
	link := NewDiLink[T0, T1](entity0, entity1)
	linkLoaded, err := link.Get(context.Background(), diLinkWrapper)
	if err != nil {
		return false, err
	}
	// load the entities
	loaded0, err := link.Get(context.Background(), link.Entity0)
	if err != nil {
		return false, err
	}
	loaded1, err := link.Get(context.Background(), link.Entity1)
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
