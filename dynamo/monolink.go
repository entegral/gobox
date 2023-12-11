package dynamo

import (
	"context"

	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"
)

// MonoLink is a generic type that can link two entities together in dynamo.
// By Default the MonoLink will establish a one-to-one relationship between the two
// entities using the primary keys. If you need to save or modify fields in the
// linked record, you will need to override this method.
type MonoLink[T0 types.Linkable] struct {
	Row // Embedding the Row type for DynamoDB requirements

	E0pk string `dynamodbav:"e0pk" json:"e0pk,omitempty"`
	E0sk string `dynamodbav:"e0sk" json:"e0sk,omitempty"`

	Entity0 T0 `dynamodbav:"-" json:"entity0,omitempty"`
}

func (m *MonoLink[T0]) LoadEntities(ctx context.Context, clients *clients.Client) (e0Loaded bool, e1Loaded bool, err error) {
	e0Loaded, err = m.LoadEntity0(ctx)
	if err != nil {
		return e0Loaded, false, err
	}
	return e0Loaded, e1Loaded, nil
}

// Link is a generic method to establish a connection between the two entities. By default
// it will establish a one-to-one relationship between the two entities using the primary keys.
// If the relation is set to OneToMany, then it will establish a one-to-many relationship
// between the two entities where Entity0 is the "one" and Entity1 is the "many".
func (m *MonoLink[T0]) Link(ctx context.Context, row types.Linkable) error {
	return m.PutLink(ctx, row)
}

// Unlink method to remove the connection between the two entities.
func (m *MonoLink[T0]) Unlink(ctx context.Context, row types.Linkable) error {
	return m.DeleteLink(ctx, row)
}
