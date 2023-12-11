package dynamo

import (
	"fmt"

	"gobox/types"

	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ErrLinkNotFound is an error type for when a new membership cannot be created.
type ErrLinkNotFound struct{}

func (e ErrLinkNotFound) Error() string {
	return "Link not found"
}

// ErrEntityNotFound is an error type for when a new membership cannot be created.
type ErrEntityNotFound[T types.Typeable] struct {
	Entity T
}

func (e ErrEntityNotFound[T]) Error() string {
	return fmt.Sprintf("%s not found", e.Entity.Type())
}

// ErrLinkTypeMismatch is an error type for when a new link should not be created
// due to a type mismatch between the dynamo row and the expected link type.
type ErrLinkTypeMismatch[T types.Typeable] struct {
	DynamoType string
	LinkType   T
}

func (e ErrLinkTypeMismatch[T]) Error() string {
	return fmt.Sprintf("Link type from dynamo '%s' does not match expected type: %s", e.DynamoType, e.LinkType.Type())
}

// ErrCouldNotValidateLink is an error type for when a link's row data cannot be validated.
type ErrCouldNotValidateLink[T types.Typeable] struct {
	LinkType T
}

func (e ErrCouldNotValidateLink[T]) Error() string {
	return fmt.Sprintf("Could not validate link type: %s", e.LinkType.Type())
}

func validateDynamoRowType[T types.Typeable](out map[string]awstypes.AttributeValue, expectedType T) error {
	v, ok := out["type"].(*awstypes.AttributeValueMemberS)
	if !ok {
		return ErrLinkTypeMismatch[T]{LinkType: expectedType}
	}
	if ok && v.Value != expectedType.Type() {
		return ErrLinkTypeMismatch[T]{LinkType: expectedType}
	}
	return nil
}
