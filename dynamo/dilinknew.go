package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// NewDiLink creates a new DiLink instance.
func NewDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) *DiLink[T0, T1] {
	link := DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}
	link.GenerateDiLinkCompositeKey()
	return &link
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
func CheckDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) (*DiLink[T0, T1], error) {
	link := NewDiLink[T0, T1](entity0, entity1)
	linkLoaded, err := GetItem(context.Background(), link)
	if err != nil {
		return link, err
	}
	// load the entities
	loaded0, err := GetItem(context.Background(), link.Entity0)
	if err != nil {
		return link, err
	}
	loaded1, err := GetItem(context.Background(), link.Entity1)
	if err != nil {
		return link, err
	}
	if loaded0 == nil || loaded0.Item == nil {
		return link, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if loaded1 == nil || loaded1.Item == nil {
		return link, ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	if linkLoaded == nil || linkLoaded.Item == nil {
		return link, ErrLinkNotFound{}
	}
	return link, nil
}
