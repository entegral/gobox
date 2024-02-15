package dynamo

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/entegral/gobox/types"
)

// EntityGSI is the name of the GSI used to contain the composite
// key of the entity at the provided index.
type EntityGSI string

// String
func (g EntityGSI) String() string {
	return string(g)
}

// GenerateMonoLinkKeys generates the composite key for the monolink.
func (m *MonoLink[T0]) GenerateMonoLinkKeys() (string, string, error) {

	m.PartitionKey = ""
	m.SortKey = ""

	var e0pk, e0sk string
	var err error

	if reflect.ValueOf(m.Entity0).IsNil() {
		if m.E0pk == "" && m.E0sk == "" {
			return "", "", errors.New("Entity0 is nil and E0pk and E0sk are empty")
		}
		e0pk = m.E0pk
		e0sk = m.E0sk
	} else {
		e0pk, e0sk, err = m.Entity0.Keys(0)
		if err != nil {
			return "", "", err
		}
	}

	linkedE0Pk, err := addKeySegment(rowType, m.Entity0.Type())
	if err != nil {
		return "", "", err
	}
	seg, err := addKeySegment(rowPk, e0pk)
	if err != nil {
		return "", "", err
	}
	linkedE0Pk += seg

	m.E0pk = linkedE0Pk
	m.E0sk = e0sk

	// Generate first part of the key using the entity0 type, pk, and sk
	seg, err = addKeySegment(entity0Type, m.Entity0.Type())
	if err != nil {
		return "", "", err
	}
	m.PartitionKey += seg
	seg, err = addKeySegment(entity0pk, e0pk)
	if err != nil {
		return "", "", err
	}
	m.PartitionKey += seg
	seg, err = addKeySegment(entity0sk, e0sk)
	if err != nil {
		return "", "", err
	}
	m.SortKey += seg
	return m.PartitionKey, m.SortKey, nil
}

// ExtractE0Keys extracts the pk and sk values for the 0th entity from the
// primary composite key.
func (m *MonoLink[T0]) ExtractE0Keys() (string, string, error) {
	if m.PartitionKey == "" || m.SortKey == "" {
		_, _, err := m.GenerateMonoLinkKeys()
		if err != nil {
			return "", "", err
		}
	}
	if m.E0pk != "" && m.E0sk != "" {
		return m.E0pk, m.E0sk, nil
	}
	pk := extractKeys(entity0pk, m.PartitionKey)
	sk := extractKeys(entity0sk, m.SortKey)
	return pk, sk, nil
}

type ErrInvalidKeySegment struct {
	label string
	value string
}

func (e ErrInvalidKeySegment) Error() string {
	return fmt.Sprintf("invalid key segment: %s(%s)", e.label, e.value)
}

func prependWithRowType(row types.Typeable, pk string) (string, error) {
	pkWithTypePrefix, err := addKeySegment(rowType, row.Type())
	if err != nil {
		return "", err
	}
	seg, err := addKeySegment(rowPk, pk)
	if err != nil {
		return "", err
	}
	pkWithTypePrefix += seg
	return pkWithTypePrefix, nil
}

func addKeySegment(label linkLabels, value string) (string, error) {
	// Check if label or value contains characters that could affect the regex
	if len(value) == 0 || strings.ContainsAny(string(label), "()") || containsObscureWhitespace(value) {
		return "", ErrInvalidKeySegment{string(label), value}
	}
	if !label.IsValidLabel() {
		return "", ErrInvalidKeySegment{string(label), value}
	}

	// Check if value matches any linkLabel
	err := label.IsValidValue(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s(%s)", label, value), nil
}

func containsObscureWhitespace(value string) bool {
	for _, r := range value {
		if unicode.IsSpace(r) && !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

// extractKeys extracts the pk and sk values from a given string.
func extractKeys(label linkLabels, str string) string {
	if !label.IsValidLabel() {
		return "invalid label"
	}

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

func (m *MonoLink[T0]) Keys(gsi int) (string, string, error) {
	// by default, we will only prefix the primary keys of both entities with "link-".
	// this will create a 1-1 relationship between the two entities.
	_, _, err := m.GenerateMonoLinkKeys()
	if err != nil {
		return "", "", err
	}

	switch gsi {
	case 0: // Primary keys
		return m.PartitionKey, m.SortKey, nil
	default:
		// Handle other GSIs or return an error
		return "", "", errors.New("invalid GSI")
	}
}
