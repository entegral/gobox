package row

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/entegral/gobox/types"
)

type ErrInvalidKeySegment struct {
	label string
	value string
}

func (e ErrInvalidKeySegment) Error() string {
	return fmt.Sprintf("invalid key segment: %s(%s)", e.label, e.value)
}

func containsObscureWhitespace(value string) bool {
	for _, r := range value {
		if unicode.IsSpace(r) && !unicode.IsPrint(r) {
			return true
		}
	}
	return false
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
