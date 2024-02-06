package row

// Link is a struct that represents a link between two entities
// It is used to represent a many-to-many relationship between two entities
// The key of the link is the combination of the keys of the two entities, applied in order
// The link can be queried on either the e0pk-e0sk-index or the e1pk-e1sk-index
// The partition keys of the primary composite key, and both GSIs, contain type
// information. This is used to make it easier to query for all links of a certain type.
type Link[T0, T1 Rowable] struct {
	Keys

	// Entity0 is the first entity in the link
	Entity0 Row[T0]

	// Entity1 is the second entity in the link
	Entity1 Row[T1]
}

// NewLink creates a new Link
func NewLink[T0, T1 Rowable](e0 Row[T0], e1 Row[T1]) Link[T0, T1] {
	// create
	return Link[T0, T1]{
		// TODO gonna jam on this later
		// Keys:    NewKeys(),
		Entity0: e0,
		Entity1: e1,
	}
}
