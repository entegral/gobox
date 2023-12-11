package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoUpdater is an interface for types so they can be used with the
// UpdateItem function.
type DynamoUpdater interface {
	Keyable
	// DynamoUpdateInput returns the UpdateItemInput for the given row.
	// While its fine to implement this method directly, it is recommended
	// to embed the row into a wrapper struct and implement this method
	// on the wrapper struct so the update behavior is appropriate for
	// the given context.
	DynamoUpdateInput(context.Context) (*dynamodb.UpdateItemInput, error)
}

// CustomDynamoMarshaller provides custom methods for marshaling and unmarshaling.
type CustomDynamoMarshaller interface {
	MarshalDynamoDB() (map[string]types.AttributeValue, error)
	UnmarshalItem(map[string]types.AttributeValue) error
}
