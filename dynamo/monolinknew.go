package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/entegral/gobox/types"
)

// CheckLink accepts both entities and attempts to load the link from dynamo.
// It does not attempt to load the entity itself, only the link.
func (link *MonoLink[T0]) CheckLink(ctx context.Context, linkWrapper types.Linkable, entity0 T0) (linkLoaded bool, err error) {
	var l *MonoLink[T0]
	if link == nil {
		l = NewMonoLink(entity0)
		link = l
	}
	loaded, err := checkMonoLink[T0](ctx, link)
	if loaded {
		err := attributevalue.UnmarshalMap(link.RowData, linkWrapper)
		return loaded, err
	}
	return false, err
}

// NewMonoLink creates a new MonoLink instance.
func NewMonoLink[T0 types.Linkable](entity0 T0) *MonoLink[T0] {
	link := MonoLink[T0]{Entity0: entity0}
	link.GenerateMonoLinkKeys()
	return &link
}

// CheckMonoLink creates a new MonoLink instance from the entities and attempts to load them from dynamo.
// If any of the entities cannot be loaded from dynamo, an error describing the missing entity will be returned.
func checkMonoLink[T0 types.Linkable](ctx context.Context, link *MonoLink[T0]) (allEntitiesExist bool, err error) {
	linkLoaded, err := link.Get(ctx, link)
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
