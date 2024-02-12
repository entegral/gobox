package row

import (
	"context"
	"fmt"

	"github.com/entegral/gobox/types"
)

// GenerateEntityKeys generates the linked entity keys using the provided entity array
func (l Link[LinkType, T]) GenerateEntityKeys(ctx context.Context) ([]Key, error) {
	var keys = make([]Key, len(l.Entities))
	for i, row := range l.Entities {
		err := row.GenerateKeys(ctx)
		if err != nil {
			return keys, err
		}
		keys[i] = Key{Pk: row.Pk, Sk: row.Sk, Index: i, IsEntity: true}
	}
	return keys, nil
}

// Link is a struct that represents a link between multiple entities
// It is used to represent a many-to-many relationship between entities
// The key of the link is the combination of the keys of the entities, applied in order
// The link can be queried on either the entityKeys-index or the entityKeys-index
// The partition keys of the primary composite key, and both GSIs, contain type
// information. This is used to make it easier to query for all links of a certain type.
type Link[LinkType types.Typeable, T Rowable] struct {
	Keys
	row LinkType

	// Entities is a slice of entities in the link
	Entities []Row[T]
}

// Type returns the type of the link
func (l Link[T0, T1]) Type() string {
	return l.row.Type()
}

func (l Link[T0, T1]) GenerateKeys(ctx context.Context) ([]Key, error) {
	keys, err := l.GenerateEntityKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error generating link entity keys for link row of type: %s. Error: %w", l.row.Type(), err)
	}

	// we need to do two things with the keys:
	// 1. Generate the primary key for this row, which is the combination of the keys of the entities, and
	// 2. Generate the keys for the entities themselves and store them in an entity GSI
	// The primary key for this row is the combination of the keys of the entities, applied in order

	var pk, sk string
	// Generate the primary key for this row
	for _, key := range keys {
		pkSeg, err := addKeySegment(rowPk, key.Pk)
		if err != nil {
			return nil, fmt.Errorf("Error adding key segment to primary key for link row of type: %s. Error: %w", l.row.Type(), err)
		}
		pk += pkSeg

		// also wrap the entity Pk with the type of the link
		key.Pk, err = prependWithRowType(l.row, key.Pk)
		if err != nil {
			return nil, fmt.Errorf("Error prepending row type to entity Pk for link row of type: %s. Error: %w", l.row.Type(), err)
		}

		skSeg, err := addKeySegment(rowSk, key.Sk)
		if err != nil {
			return nil, fmt.Errorf("Error adding key segment to primary key for link row of type: %s. Error: %w", l.row.Type(), err)
		}
		sk += skSeg
	}

	// prefix the primary key with the type of the link
	l.Pk, err = prependWithRowType(l.row, l.Pk)
	return keys, nil
}

// NewLink creates a new Link
func NewLink[LinkType types.Typeable, T Rowable](link LinkType, entities []Row[T]) Link[LinkType, T] {
	return Link[LinkType, T]{
		row:      link,
		Entities: entities,
	}
}
