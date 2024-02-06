package row

import (
	"context"
	"os"

	"github.com/entegral/gobox/clients"
)

type Table struct {
	tableName string
	client    *clients.Client
}

func (t *Table) TableName() string {
	// If the table name is not set, use the value provided from TABLE_NAME
	if t.tableName == "" {
		return os.Getenv("TABLE_NAME")
	}
	return t.tableName
}

func (t *Table) SetTableName(tableName string) {
	t.tableName = tableName
}

func (t *Table) SetClient(client *clients.Client) {
	t.client = client
}

func (t *Table) GetClient(ctx context.Context) *clients.Client {
	// If the client is not set, use the default client from the clients package
	if t.client != nil {
		return t.client
	}
	return clients.GetDefaultClient(ctx)
}
