package dynamo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RowExample struct {
	Row
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (re RowExample) Type() string {
	return "RowExample"
}

// Keys returns the partition key and sort key for the row
func (re *RowExample) Keys(gsi int) (string, string, error) {
	// For this example, assuming GUID is the partition key and Email is the sort key.
	// Additional logic can be added to handle different GSIs if necessary.
	switch gsi {
	case 0: // Primary keys
		return "partitionKey", "sortKey", nil
	default:
		// Handle other GSIs or return an error
		return "", "", nil
	}
}

type MonoLinkExample struct {
	*MonoLink[*RowExample]
}

func (mle *MonoLinkExample) Type() string {
	return "MonoLinkExample"
}

func TestMonoLinkKeyGen(t *testing.T) {
	// Create a new row
	row := &RowExample{
		Title:       "Test Title",
		Description: "Test Description",
	}

	// Create a new monolink
	monoLink := &MonoLinkExample{
		MonoLink: NewMonoLink(row),
	}

	// Generate the composite key
	linkPk, linkSk, _ := monoLink.GenerateMonoLinkKeys()

	// Verify the keys
	assert.Equal(t, "/e0Type(RowExample)/e0pk(partitionKey)", linkPk)
	assert.Equal(t, "/e0sk(sortKey)", linkSk)

	// Extract the keys
	pk, sk, _ := monoLink.ExtractE0Keys()
	assert.Equal(t, "/rowType(RowExample)/rowPk(partitionKey)", pk)
	assert.Equal(t, "sortKey", sk)
}

func TestExtracKeys(t *testing.T) {
	const key1 = "/e0Type(RowExample)/e0pk(partitionKey)"
	const key2 = "/e0sk(sortKey)"
	const key3 = "/rowType(RowExample)/rowPk(partitionKey)"
	t.Run("entity0", func(t *testing.T) {
		t.Run("type", func(t *testing.T) {
			typeTest := extractKeys(entity0Type, key1)
			fmt.Println(typeTest)
			assert.Equal(t, "RowExample", typeTest)
		})
		t.Run("pk", func(t *testing.T) {
			pkTest := extractKeys(entity0pk, key1)
			fmt.Println(pkTest)
			assert.Equal(t, "partitionKey", pkTest)
		})
		t.Run("sk", func(t *testing.T) {
			skTest := extractKeys(entity0sk, key2)
			fmt.Println(skTest)
			assert.Equal(t, "sortKey", skTest)
		})
	})
	t.Run("row", func(t *testing.T) {
		t.Run("type", func(t *testing.T) {
			rowTypeTest := extractKeys(rowType, key3)
			fmt.Println(rowTypeTest)
			assert.Equal(t, "RowExample", rowTypeTest)
		})
		t.Run("pk", func(t *testing.T) {
			rowPkTest := extractKeys(rowPk, key3)
			fmt.Println(rowPkTest)
			assert.Equal(t, "partitionKey", rowPkTest)
		})
	})
}

func TestAddKeySegment(t *testing.T) {
	tests := []struct {
		name    string
		label   linkLabels
		value   string
		want    string
		wantErr bool
	}{
		{
			name:  "Test with non-empty label",
			label: "rowType",
			value: "testValue",
			want:  "/rowType(testValue)",
		},
		{
			name:    "Test with unknown label",
			label:   "test/Label",
			value:   "testValue",
			wantErr: true,
		},
		{
			name:  "Test with non-alphanumeric value",
			label: "rowType",
			value: "test/Value",
			want:  "/rowType(test/Value)",
		},
		{
			name:    "Test with empty value",
			label:   "rowType",
			value:   "",
			wantErr: true,
		},
		{
			name:    "Test with both label and value empty",
			label:   "",
			value:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test with value containing newline character",
			label:   "testLabel",
			value:   "test\nValue",
			wantErr: true,
		},
		{
			name:    "Test with value matching a linkLabel",
			label:   "testLabel",
			value:   "e0Type",
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(i)
			got, err := addKeySegment(tt.label, tt.value)
			if (err == nil) && tt.wantErr {
				t.Errorf("addKeySegment() error = %v, wantErr %v, for %v", err, tt.wantErr, tt.name)
				return
			}
			if got != tt.want {
				t.Errorf("addKeySegment() = %v, want %v", got, tt.want)
			}
		})
	}
}
