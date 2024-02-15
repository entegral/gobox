package dynamo

import (
	"context"
	"os"

	"github.com/entegral/gobox/clients"
)

type Table struct {
	Client    *clients.Client
	Tablename string
}

func NewTable(tablename string) Table {
	return Table{Tablename: tablename}
}

// TableName returns the name of the DynamoDB table.
// By default, this is the value of the TABLENAME environment variable.
// If you need to override this, implement this method on the parent type.
func (t *Table) TableName(ctx context.Context) string {
	if t.Tablename != "" {
		return t.Tablename
	}
	tn := os.Getenv("TABLENAME")
	if tn == "" {
		panic("TABLENAME environment variable not set")
	}
	return tn
}

func (t *Table) SetTableName(tablename string) {
	t.Tablename = tablename
}

// SetClient sets the client to use for DynamoDB operations.
func (t *Table) SetClient(client *clients.Client) {
	t.Client = client
}
