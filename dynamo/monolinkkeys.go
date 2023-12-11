package dynamo

import (
	"fmt"
	"regexp"
)

// EntityGSI is the name of the GSI used to contain the composite
// key of the entity at the provided index.
type EntityGSI string

// String
func (g EntityGSI) String() string {
	return string(g)
}

// GenerateMonoLinkCompositeKey generates the composite key for the monolink.
func (m *MonoLink[T0]) GenerateMonoLinkCompositeKey() (string, string) {
	m.Pk = ""
	m.Sk = ""
	e0pk, e0sk := m.Entity0.Keys(0)
	m.E0pk = e0pk
	m.E0sk = e0sk

	// Generate first part of the key using the entity0 type, pk, and sk
	m.Pk += addKeySegment(entity0Type, m.Entity0.Type())
	m.Pk += addKeySegment(entity0pk, e0pk)
	m.Sk += addKeySegment(entity0sk, e0sk)
	return m.Pk, m.Sk
}

// ExtractE0Keys extracts the pk and sk values for the 0th entity from the
// primary composite key.
func (m *MonoLink[T0]) ExtractE0Keys() (string, string) {
	if m.Pk == "" || m.Sk == "" {
		m.GenerateMonoLinkCompositeKey()
	}
	if m.E0pk != "" && m.E0sk != "" {
		return m.E0pk, m.E0sk
	}
	pk := extractKeys(entity0pk, m.Pk)
	sk := extractKeys(entity0sk, m.Sk)
	linkedPk := addKeySegment(rowType, m.Type())
	linkedPk += addKeySegment(entity0pk, pk)
	return linkedPk, sk
}

func addKeySegment(label linkLabels, value string) string {
	if label == "" {
		return value
	}
	return fmt.Sprintf("/%s(%s)", label, value)
}

// extractKeys extracts the pk and sk values from a given string.
func extractKeys(label linkLabels, str string) string {
	// Define regular expressions for pk and sk

	// regexFormat - where %d is the entity number and %s either Pk or Sk
	regexFormat := `(?m)%s\(([^)]+)\)`

	regex := regexp.MustCompile(fmt.Sprintf(regexFormat, label))

	// Find pk and sk
	pkMatches := regex.FindStringSubmatch(str)
	if len(pkMatches) == 2 {
		return pkMatches[1]
	}
	return "nothing found"
}

func (m *MonoLink[T0]) Keys(gsi int) (string, string) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _ = m.GenerateMonoLinkCompositeKey()

	switch gsi {
	case 0: // Primary keys
		return m.Pk, m.Sk
	default:
		// Handle other GSIs or return an error
		return "", ""
	}
}
