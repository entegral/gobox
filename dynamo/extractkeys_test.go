package dynamo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKeys(t *testing.T) {
	str := "/rowType(PinkSlip)/rowPk(/e0Type(user)/e0pk(testDiLinkEmail@gmail.com)/e1Type(car)/e1pk(TestDiLinkMake2-TestDiLinkModel2))"

	t.Run("Extract pk and sk for e0Type", func(t *testing.T) {
		key := extractKeys(entity0Type, str)
		assert.Equal(t, "user", key)
	})

	t.Run("Extract pk and sk for e1Type", func(t *testing.T) {
		key := extractKeys(entity1Type, str)
		assert.Equal(t, "car", key)
	})

	t.Run("Extract pk and sk for rowType", func(t *testing.T) {
		key := extractKeys(rowType, str)
		assert.Equal(t, "PinkSlip", key)
	})

	t.Run("Extract pk and sk for non-existent label", func(t *testing.T) {
		key := extractKeys("", str)
		assert.Equal(t, "invalid label", key)
	})
}
