package dynamo

func (m *TriLink[T0, T1, T2]) GenerateTriLinkCompositeKey() (string, string, error) {
	// generate keys for the 0th entity
	_, _, err := m.GenerateDiLinkKeys()
	if err != nil {
		return "", "", err
	}

	// Generate third part of the key using the entity2 type, pk, and sk
	// to ensure uniqueness of the key
	e2pk, e2sk, err := m.Entity2.Keys(0)
	if err != nil {
		return "", "", err
	}

	linkedE2Pk, err := addKeySegment(rowType, m.Entity2.Type())
	if err != nil {
		return "", "", err
	}
	seg, err := addKeySegment(rowPk, e2pk)
	if err != nil {
		return "", "", err
	}
	linkedE2Pk += seg

	m.E2pk = linkedE2Pk
	m.E2sk = e2sk

	seg, err = addKeySegment(entity2Type, m.Entity2.Type())
	if err != nil {
		return "", "", err
	}
	m.Pk += seg
	seg, err = addKeySegment(entity2pk, e2pk)
	if err != nil {
		return "", "", err
	}
	m.Pk += seg
	seg, err = addKeySegment(entity2sk, e2sk)
	if err != nil {
		return "", "", err
	}
	m.Sk += seg
	return m.Pk, m.Sk, nil
}

func (m *TriLink[T0, T1, T2]) ExtractE2Keys() (string, string, error) {
	if m.Pk == "" || m.Sk == "" {
		_, _, err := m.GenerateTriLinkCompositeKey()
		if err != nil {
			return "", "", err
		}
	}
	if m.E2pk != "" && m.E2sk != "" {
		return m.E2pk, m.E2sk, nil
	}
	pk2 := extractKeys(entity2pk, m.Pk)
	sk2 := extractKeys(entity2sk, m.Sk)

	return pk2, sk2, nil
}

func (m *TriLink[T0, T1, T2]) Keys(gsi int) (string, string, error) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _, err := m.GenerateTriLinkCompositeKey()
	if err != nil {
		return "", "", err
	}

	switch gsi {
	case 0: // Primary keys
		return m.Pk, m.Sk, nil
	case 1: // GSI 1
		return m.Pk1, m.Sk1, nil
	default:
		// Handle other GSIs or return an error
		return "", "", ErrInvalidGSI{GSI: gsi}
	}
}
