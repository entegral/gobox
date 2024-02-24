package dynamo

// KeyableTimeCapsule implements the Keyable interface and will use the Location and Name fields to generate the pk and sk
type KeyableTimeCapsule struct {
	Row
	Name     string `dynamo:"name"`     // Name of the TimeCapsule
	Location string `dynamo:"location"` // Location of the TimeCapsule
}

func (tc *KeyableTimeCapsule) Keys(gsi int) (string, string, error) {
	// our implementation will use the Name and Location fields to generate the pk and sk
	tc.PartitionKey = tc.Location
	tc.SortKey = tc.Name
	return tc.PartitionKey, tc.SortKey, nil
}
