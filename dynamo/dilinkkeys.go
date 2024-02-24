package dynamo

import (
	"errors"
	"reflect"
)

func (m *DiLink[T0, T1]) GenerateDiLinkKeys() (string, string, error) {

	// generate keys for the 0th entity
	_, _, err := m.GenerateMonoLinkKeys()
	if err != nil {
		return "", "", err
	}

	var e1pk, e1sk string

	if reflect.ValueOf(m.Entity1).IsNil() {
		if m.E1pk == "" && m.E1sk == "" {
			return "", "", errors.New("Entity1 is nil and E1pk and E1sk are empty")
		}
		e1pk = extractKeys(rowPk, m.E1pk)
		e1sk = m.E1sk
	} else {
		// Generate second part of the key using the entity1 type, pk, and sk
		// to ensure uniqueness of the key
		e1pk, e1sk, err = m.Entity1.Keys(0)
		if err != nil {
			return "", "", err
		}
	}

	linkedE1Pk, errPk := addKeySegment(rowType, m.Entity1.Type())
	if errPk != nil {
		return "", "", errPk
	}
	seg, errPk2 := addKeySegment(rowPk, e1pk)
	if errPk2 != nil {
		return "", "", errPk2
	}
	linkedE1Pk += seg

	m.E1pk = linkedE1Pk
	m.E1sk = e1sk

	seg, errPk = addKeySegment(entity1Type, m.Entity1.Type())
	if errPk != nil {
		return "", "", errPk
	}
	m.PartitionKey += seg
	seg, errPk2 = addKeySegment(entity1pk, e1pk)
	if errPk2 != nil {
		return "", "", errPk2
	}
	m.PartitionKey += seg
	seg, errSk := addKeySegment(entity1sk, e1sk)
	if errSk != nil {
		return "", "", errSk
	}
	m.SortKey += seg
	return m.PartitionKey, m.SortKey, nil
}

func (m *DiLink[T0, T1]) ExtractE1Keys() (string, string, error) {
	if m.PartitionKey == "" || m.SortKey == "" {
		_, _, err := m.GenerateDiLinkKeys()
		if err != nil {
			return "", "", err
		}
	}
	if m.E1pk != "" && m.E1sk != "" {
		return m.E1pk, m.E1sk, nil
	}
	pk1 := extractKeys(entity1pk, m.PartitionKey)
	sk1 := extractKeys(entity1sk, m.SortKey)
	return pk1, sk1, nil
}

func (m *DiLink[T0, T1]) Keys(gsi int) (string, string, error) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _, err := m.GenerateDiLinkKeys()
	if err != nil {
		return "", "", err
	}

	switch gsi {
	case 0: // Primary keys
		return m.PartitionKey, m.SortKey, nil
	case 1: // GSI 1
		return *m.Pk1, *m.Sk1, nil
	default:
		// Handle other GSIs or return an error
		return "", "", ErrInvalidGSI{GSI: gsi}
	}
}
