package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// TriLink is a generic type that can link two entities together in dynamo.
// By Default the TriLink will establish a one-to-one relationship between the two
// entities using the primary keys. If you need to save or modify fields in the
// linked record, you will need to override this method.
type TriLink[T0, T1, T2 types.Linkable] struct {
	DiLink[T0, T1] // Embedding the DiLink type for DynamoDB requirements

	E2pk string `dynamodbav:"e2pk" json:"e1pk,omitempty"`
	E2sk string `dynamodbav:"e2sk" json:"e1sk,omitempty"`

	Entity2 T2 `dynamodbav:"-" json:"entity1,omitempty"`
}

func (m *TriLink[T0, T1, T2]) LoadEntities(ctx context.Context) (e0Loaded, e1Loaded, e2Loaded bool, err error) {
	e0Loaded, err = m.LoadEntity0(ctx)
	if err != nil {
		return e0Loaded, false, false, err
	}
	e1Loaded, err = m.LoadEntity1(ctx)
	if err != nil {
		return e0Loaded, e1Loaded, false, err
	}
	e2Loaded, err = m.LoadEntity2(ctx)
	if err != nil {
		return e0Loaded, e1Loaded, e2Loaded, err
	}
	return e0Loaded, e1Loaded, e2Loaded, nil
}

// Link is a generic method to establish a connection between the two entities. By default
// it will establish a one-to-one relationship between the two entities using the primary keys.
// If the relation is set to OneToMany, then it will establish a one-to-many relationship
// between the two entities where Entity0 is the "one" and Entity1 is the "many".
func (m *TriLink[T0, T1, T2]) Link(ctx context.Context, row types.Linkable) error {
	return m.Put(ctx, row)
}

// Unlink method to remove the connection between the two entities.
func (m *TriLink[T0, T1, T2]) Unlink(ctx context.Context, row types.Linkable) error {
	return m.Delete(ctx, row)
}
