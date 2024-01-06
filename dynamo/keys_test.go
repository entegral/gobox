package dynamo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type a struct {
	Row
}

func (a a) Keys(gsi int) (string, string, error) {
	return "partionKeyA", "sortKeyA", nil
}

func (a a) Type() string {
	return "typeA"
}

type b struct {
	Row
	value int
}

func (b b) Type() string {
	if b.value == 0 {
		return "typeB"
	}
	return "typeB" + fmt.Sprintf("%d", b.value)
}

func (b b) Keys(gsi int) (string, string, error) {
	return "partionKeyB", "sortKeyB", nil
}

type testType struct {
	DiLink[a, b]
}

func TestGenerateLinkKeys(t *testing.T) {
	t.Run("Generate Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkKeys()
		assert.Equal(t, "/e0Type(typeA)/e0pk(partionKeyA)/e1Type(typeB)/e1pk(partionKeyB)", testType.Pk)
		assert.Equal(t, "/e0sk(sortKeyA)/e1sk(sortKeyB)", testType.Sk)
	})
	t.Run("Extract Entity0 Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkKeys()
		pk, sk, _ := testType.ExtractE0Keys()
		assert.Equal(t, "/rowType(typeA)/rowPk(partionKeyA)", pk)
		assert.Equal(t, "sortKeyA", sk)
	})
	t.Run("Extract Entity1 Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkKeys()
		pk, sk := testType.ExtractE1Keys()
		assert.Equal(t, "/rowType(typeB)/rowPk(partionKeyB)", pk)
		assert.Equal(t, "sortKeyB", sk)
	})
}
