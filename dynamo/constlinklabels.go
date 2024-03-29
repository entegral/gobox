package dynamo

import (
	"errors"
)

type linkLabels string

func (l linkLabels) String() string {
	return string(l)
}

var validLabels = []linkLabels{
	pk,
	sk,
	entity0pk,
	entity0sk,
	entity0Type,
	rowType,
	rowPk,
	rowSk,
	entity2pk,
	entity2sk,
	entity2Type,
	entity1pk,
	entity1sk,
	entity1Type,
}

func (ll linkLabels) IsValidLabel() bool {
	if ll == "" {
		return false
	}
	for _, label := range validLabels {
		if label.String() == ll.String() {
			return true
		}
	}
	return false
}

func (ll linkLabels) IsValidValue(value string) error {
	for _, label := range validLabels {
		if label.String() == value {
			return errors.New("value must not match any linkLabel")
		}
	}
	return nil
}

func (ll linkLabels) Values() []string {
	values := make([]string, len(ll))
	for i, value := range ll {
		values[i] = string(value)
	}
	return values
}

const (
	// Entity0GSI is the name of the GSI used to contain the primary composite
	// key of the 0th entity.
	Entity0GSI  EntityGSI  = "e0pk-e0sk-index"
	pk          linkLabels = "pk"
	sk          linkLabels = "sk"
	entity0pk   linkLabels = "e0pk"
	entity0sk   linkLabels = "e0sk"
	entity0Type linkLabels = "e0Type"
	rowType     linkLabels = "rowType"
	rowPk       linkLabels = "rowPk"
	rowSk       linkLabels = "rowSk"
)

const (
	// Entity1GSI is the name of the GSI used to contain the primary composite
	// key of the 1st entity.
	Entity1GSI EntityGSI = "e1pk-e1sk-index"

	entity1pk   linkLabels = "e1pk"
	entity1sk   linkLabels = "e1sk"
	entity1Type linkLabels = "e1Type"
)

const (
	// Entity2GSI is the name of the GSI used to contain the primary composite
	// key of the 2nd entity.
	Entity2GSI EntityGSI = "e2pk-e2sk-index"

	entity2pk   linkLabels = "e2pk"
	entity2sk   linkLabels = "e2sk"
	entity2Type linkLabels = "e2Type"
)
