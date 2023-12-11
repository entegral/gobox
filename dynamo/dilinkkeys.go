package dynamo

func (m *DiLink[T0, T1]) GenerateDiLinkCompositeKey() (string, string) {
	// generate keys for the 0th entity
	m.GenerateMonoLinkCompositeKey()

	// Generate second part of the key using the entity1 type, pk, and sk
	// to ensure uniqueness of the key
	e1pk, e1sk := m.Entity1.Keys(0)

	linkedE1Pk := addKeySegment(rowType, m.Type())
	linkedE1Pk += addKeySegment(entity1pk, e1pk)

	m.E1pk = linkedE1Pk
	m.E1sk = e1sk

	m.Pk += addKeySegment(entity1Type, m.Entity1.Type())
	m.Pk += addKeySegment(entity1pk, e1pk)
	m.Sk += addKeySegment(entity1sk, e1sk)
	return m.Pk, m.Sk
}

func (m *DiLink[T0, T1]) ExtractE1Keys() (string, string) {
	if m.Pk == "" || m.Sk == "" {
		m.GenerateDiLinkCompositeKey()
	}
	if m.E1pk != "" && m.E1sk != "" {
		return m.E1pk, m.E1sk
	}
	pk1 := extractKeys(entity1pk, m.Pk)
	sk1 := extractKeys(entity1sk, m.Sk)

	return pk1, sk1
}

func (m *DiLink[T0, T1]) Keys(gsi int) (string, string) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _ = m.GenerateDiLinkCompositeKey()

	switch gsi {
	case 0: // Primary keys
		return m.Pk, m.Sk
	case 1: // GSI 1
		return m.Pk1, m.Sk1
	default:
		// Handle other GSIs or return an error
		return "", ""
	}
}
