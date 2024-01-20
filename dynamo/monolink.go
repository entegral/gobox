package dynamo

import (
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

// Type returns the type of the record.
func (r *MonoLink[T0]) Type() string {
	if r.UnmarshalledType == "" {
		return "MonoLink"
	}
	return r.UnmarshalledType
}
