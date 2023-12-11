package dynamo

func (m *TriLink[T0, T1, T2]) GenerateTriLinkCompositeKey() (string, string) {
	// generate keys for the 0th entity
	m.GenerateDiLinkCompositeKey()

	// Generate second part of the key using the entity1 type, pk, and sk
	// to ensure uniqueness of the key
	e2pk, e2sk := m.Entity2.Keys(0)
	m.E2pk = e2pk
	m.E2sk = e2sk
	m.Pk += addKeySegment(entity2Type, m.Entity2.Type())
	m.Pk += addKeySegment(entity2pk, e2pk)
	m.Sk += addKeySegment(entity2sk, e2sk)
	return m.Pk, m.Sk
}

func (m *TriLink[T0, T1, T2]) ExtractE2Keys() (string, string) {
	if m.Pk == "" || m.Sk == "" {
		m.GenerateTriLinkCompositeKey()
	}
	if m.E1pk != "" && m.E1sk != "" {
		return m.E2pk, m.E2sk
	}
	pk1 := extractKeys(entity2pk, m.Pk)
	sk1 := extractKeys(entity2sk, m.Sk)

	return pk1, sk1
}

func (m *TriLink[T0, T1, T2]) Keys(gsi int) (string, string) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _ = m.GenerateTriLinkCompositeKey()

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
