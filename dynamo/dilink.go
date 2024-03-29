package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"
)

// DiLink is a generic type that can link two entities together in dynamo.
// By Default the DiLink will establish a one-to-one relationship between the two
// entities using the primary keys. If you need to save or modify fields in the
// linked record, you will need to override this method.
type DiLink[T0 types.Linkable, T1 types.Linkable] struct {
	MonoLink[T0] // Embedding the MonoLink type for DynamoDB requirements

	E1pk string `dynamodbav:"e1pk" json:"e1pk,omitempty"`
	E1sk string `dynamodbav:"e1sk" json:"e1sk,omitempty"`

	Entity1 T1 `dynamodbav:"-" json:"entity1,omitempty"`
}

func (m *DiLink[T0, T1]) LoadEntities(ctx context.Context) (e0Loaded bool, e1Loaded bool, err error) {
	e0Loaded, err = m.LoadEntity0(ctx)
	if err != nil {
		return e0Loaded, false, err
	}
	e1Loaded, err = m.LoadEntity1(ctx)
	if err != nil {
		return e0Loaded, e1Loaded, err
	}
	return e0Loaded, e1Loaded, nil
}

// Type returns the type of the record.
func (r *DiLink[T0, T1]) Type() string {
	if r.UnmarshalledType == "" {
		return "dilink"
	}
	return r.UnmarshalledType
}

// Link is a generic method to establish a connection between the two entities.
// Any two entities that embed the Row type can be linked together while maintaining
// primary key entropy equal to the sum of the two entities.
func (m *DiLink[T0, T1]) Link(ctx context.Context, row types.Linkable) error {
	return m.Put(ctx, row)
}

// Unlink method removes the connection between the two entities by deleting the link record.
func (m *DiLink[T0, T1]) Unlink(ctx context.Context, row types.Linkable) error {
	return m.Delete(ctx, row)
}
