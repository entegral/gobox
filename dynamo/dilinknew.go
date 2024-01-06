package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entities themselves, only the link.
func (link *DiLink[T0, T1]) CheckLink(ctx context.Context, linkWrapper types.Linkable, entity0 T0, entity1 T1) (loaded bool, err error) {
	var l *DiLink[T0, T1]
	l = NewDiLink(entity0, entity1)
	link = l
	loaded, err = checkDiLink[T0, T1](ctx, link)
	if loaded {
		err := attributevalue.UnmarshalMap(link.RowData, linkWrapper)
		return loaded, err
	}
	return false, err
}

// NewDiLink creates a new DiLink instance.
func NewDiLink[T0, T1 types.Linkable](entity0 T0, entity1 T1) *DiLink[T0, T1] {
	link := DiLink[T0, T1]{MonoLink: MonoLink[T0]{Entity0: entity0}, Entity1: entity1}
	link.GenerateDiLinkKeys()
	return &link
}

// CheckDiLink creates a new DiLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func checkDiLink[T0, T1 types.Linkable](ctx context.Context, link *DiLink[T0, T1]) (allEntitiesExist bool, err error) {
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
	linkLoaded, err := link.Get(ctx, link)
	if err != nil {
		return false, err
	}
	if !linkLoaded {
		return true, ErrLinkNotFound{}
	}
	return true, nil
}
