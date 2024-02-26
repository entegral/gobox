package dynamo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_Query(t *testing.T) {

	// Create a new Table
	table := NewTable("arctica")

	// Create a new GSI
	gsi := ItemKey{
		PkPrefix: "pk2",
		SkPrefix: "sk2",
		PkValue:  "partitionKey",
		SkValue:  "sortKey",
	}

	// Execute the Query method
	output, err := table.Query(context.Background(), gsi)

	// Assert that there was no error and the output is not nil
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, table.Queries, 1)
	assert.Len(t, table.Queries[0].Items, 0)
}
