package row

import "fmt"

type GSI string // GSI is a global secondary index

func (g GSI) String() string {
	return string(g)
}

const (
	PK0SK0Index GSI = "pk0-sk0-index"
	PK1SK1Index GSI = "pk1-sk1-index"
	PK2SK2Index GSI = "pk2-sk2-index"
	PK3SK3Index GSI = "pk3-sk3-index"
	PK4SK4Index GSI = "pk4-sk4-index"
	PK5SK5Index GSI = "pk5-sk5-index"
	PK6SK6Index GSI = "pk6-sk6-index"
)

const (
	// Entity0GSI is the name of the GSI used to contain the primary composite
	Entity0GSI GSI = "e0pk-e0sk-index"

	// Entity1GSI is the name of the GSI used to contain the primary composite
	// key of the 1st entity.
	Entity1GSI GSI = "e1pk-e1sk-index"

	// Entity2GSI is the name of the GSI used to contain the primary composite
	// key of the 2nd entity.
	Entity2GSI GSI = "e2pk-e2sk-index"
)

func GetIndexName(key Key) string {
	if !key.IsEntity {
		return fmt.Sprintf("pk%d-sk%d-index", key.Index, key.Index)
	}
	switch key.Index {
	case 0:
		return Entity0GSI.String()
	case 1:
		return Entity1GSI.String()
	case 2:
		return Entity2GSI.String()
	default:
		return "unknown index"
	}
}
