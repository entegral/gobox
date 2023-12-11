package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// NewTriLink creates a new TriLink instance.
func NewTriLink[T0, T1, T2 types.Linkable](entity0 T0, entity1 T1, entity2 T2) *TriLink[T0, T1, T2] {
	link := TriLink[T0, T1, T2]{DiLink: DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}, Entity2: entity2}
	link.GenerateTriLinkCompositeKey()
	return &link
}

// CheckTriLink creates a new TriLink instance and attempts to load the entities.
func CheckTriLink[T0, T1, T2 types.Linkable](entity0 T0, entity1 T1, entity2 T2) (*TriLink[T0, T1, T2], error) {
	link := NewTriLink[T0, T1, T2](entity0, entity1, entity2)
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
	loaded2, err := GetItem(context.Background(), link.Entity2)
	if err != nil {
		return link, err
	}
	if loaded0 == nil || loaded0.Item == nil {
		return link, ErrEntityNotFound[T0]{Entity: link.Entity0}
	}
	if loaded1 == nil || loaded1.Item == nil {
		return link, ErrEntityNotFound[T1]{Entity: link.Entity1}
	}
	if loaded2 == nil || loaded2.Item == nil {
		return link, ErrEntityNotFound[T2]{Entity: link.Entity2}
	}
	if linkLoaded == nil || linkLoaded.Item == nil {
		return link, ErrLinkNotFound{}
	}
	return link, nil
}
