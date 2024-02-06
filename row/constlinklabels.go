package row

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
	// use switch case instead of for loop
	switch ll {
	case pk, sk, entity0pk, entity0sk, entity0Type, rowType, rowPk, rowSk, entity2pk, entity2sk, entity2Type, entity1pk, entity1sk, entity1Type:
		return true
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
	// key of the 0th entity.
	pk          linkLabels = "pk"
	sk          linkLabels = "sk"
	entity0pk   linkLabels = "e0pk"
	entity0sk   linkLabels = "e0sk"
	entity0Type linkLabels = "e0#"
	rowType     linkLabels = "r#"
	rowPk       linkLabels = "rPk"
	rowSk       linkLabels = "rSk"
)

const (
	entity1pk   linkLabels = "e1pk"
	entity1sk   linkLabels = "e1sk"
	entity1Type linkLabels = "e1#"
)

const (
	entity2pk   linkLabels = "e2pk"
	entity2sk   linkLabels = "e2sk"
	entity2Type linkLabels = "e2#"
)
