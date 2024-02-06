package row

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/types"
)

type Rowable interface {
	types.Keyable
	types.Typeable
}

type Row[T Rowable] struct {
	// Table configuration, marshaled to json but not to dynamo
	Table `dynamodbav:"-" json:"tablename,omitempty"`

	// The keys for the row
	Keys

	// The object that is being stored or retrieved
	object T
}

func (r *Row[T]) Type() string {
	return r.object.Type()
}

func (r *Row[T]) unmarshalMap(m map[string]awstypes.AttributeValue) error {
	// Create a new map to hold the non-key values
	err := attributevalue.UnmarshalMap(m, &r.Keys)
	if err != nil {
		return err
	}

	// Unmarshal the non-key values into the Object
	err = attributevalue.UnmarshalMap(m, &r.object)
	return err
}

// Object returns the object that is being stored or retrieved
func (r *Row[T]) Object() T {
	return r.object
}

func NewRow[T Rowable](object T) Row[T] {
	return Row[T]{object: object}
}
