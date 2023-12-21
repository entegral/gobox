package types

import (
	"context"

	"github.com/entegral/gobox/clients"
)

type Tableable interface {
	TableName(ctx context.Context) string
}

func CheckTableable(ctx context.Context, item interface{}) string {
	var tn string
	// check if row implements the Tableable interface
	if tableable, ok := item.(Tableable); ok {
		// if so, use the provided tablename
		tn = tableable.TableName(ctx)
	} else {
		// otherwise, use the default tablename
		tn = clients.TableName(ctx)
	}
	return tn
}
