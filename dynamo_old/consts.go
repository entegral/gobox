package dynamo

// GSIName is a list of consts that represent the names of the GSIs
// that should be available on a dynamodb table.
type GSIName string

const (
	// GSI1 is the name of the first GSI
	GSI1 GSIName = "pk1-sk1-index"
	// GSI2 is the name of the second GSI
	GSI2 GSIName = "pk2-sk2-index"
	// GSI3 is the name of the third GSI
	GSI3 GSIName = "pk3-sk3-index"
	// GSI4 is the name of the fourth GSI
	GSI4 GSIName = "pk4-sk4-index"
	// GSI5 is the name of the fifth GSI
	GSI5 GSIName = "pk5-sk5-index"
	// GSI6 is the name of the sixth GSI
	GSI6 GSIName = "pk6-sk6-index"
)

// String returns the name of the GSI
func (g GSIName) String() string {
	return string(g)
}
