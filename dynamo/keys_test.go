package dynamo

import (
	"fmt"
	"testing"
)

type a struct {
	Row
}

func (a a) Keys(gsi int) (string, string) {
	return "partionKeyA", "sortKeyA"
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

func (b b) Keys(gsi int) (string, string) {
	return "partionKeyB", "sortKeyB"
}

type testType struct {
	DiLink[a, b]
}

func TestGenerateLinkCompositeKey(t *testing.T) {
	t.Run("Generate Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkCompositeKey()
		if testType.Pk != "/e0Type(typeA)/e0Pk(partionKeyA)/e1Type(typeB)/e1Pk(partionKeyB)" {
			t.Errorf("Expected Pk to be /e0Type(typeA)/e0Pk(partionKeyA)/e1Type(typeB)/e1Pk(partionKeyB), got %s", testType.Pk)
		}
		if testType.Sk != "/e0Sk(sortKeyA)/e1Sk(sortKeyB)" {
			t.Errorf("Expected Sk to be /e0Sk(sortKeyA)/e1Sk(sortKeyB), got %s", testType.Sk)
		}
	})
	t.Run("Extract Entity0 Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkCompositeKey()
		pk, sk := testType.ExtractE0Keys()
		if pk != "partionKeyA" {
			fmt.Println("testType.Pk", testType.Pk)
			t.Errorf("Expected pk to be partionKeyA, got %s", pk)
		}
		if sk != "sortKeyA" {
			fmt.Println("testType.Sk", testType.Sk)
			t.Errorf("Expected sk to be sortKeyA, got %s", sk)
		}
	})
	t.Run("Extract Entity1 Pk and Sk", func(t *testing.T) {
		a := a{}
		b := b{}
		testType := testType{}
		testType.Entity0 = a
		testType.Entity1 = b
		testType.GenerateDiLinkCompositeKey()
		pk, sk := testType.ExtractE1Keys()
		if pk != "partionKeyB" {
			fmt.Println("testType.Pk", testType.Pk)
			t.Errorf("Expected pk to be partionKeyA, got %s", pk)
		}
		if sk != "sortKeyB" {
			fmt.Println("testType.Sk", testType.Sk)
			t.Errorf("Expected sk to be sortKeyA, got %s", sk)
		}
	})
}
