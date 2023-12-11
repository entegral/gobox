package types

// Keyable is an interface for types so they can be used with the
// GetItem, PutItem, and DeleteItem functions.
type Keyable interface {

	// Keys returns the partition key and sort key for the given GSI.
	// If the GSI is 0, then the primary composite key is returned/assumed.
	// If the GSI is 1, then the composite key for the pk1-sk1-index is assumed.
	// When implementing this method, you should return the appropriate
	// partition key and sort key for the given GSI, however, you should
	// also ensure any other GSI fields that rely on struct fields are
	// populated as well.
	Keys(gsi int) (partitionKey, sortKey string)
}
