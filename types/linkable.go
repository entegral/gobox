package types

import "context"

type Typeable interface {
	// Type returns the type of the row.
	// By default, this is the name of the struct, but it can be overridden
	// by implementing this method.
	Type() string
}

type Tableable interface {
	TableName(ctx context.Context) string
}

// Linkable is an interface that can be embedded in a struct to indicate that it
// can be linked to another entity. This interface is used by the DiLink type and
// requires that the struct implement the Keyable and Typeable interfaces.
type Linkable interface {
	Keyable
	Typeable
	Tableable
}
