package dynamo

import (
	"github.com/dgryski/trifles/uuid"
	"github.com/entegral/gobox/keys"
)

// Row is a sample Keyable implementation. It is not intended to be used
// by itself, but rather to be embedded into other types. After embedding,
// you should implement the TableName and Keys methods on the parent type.
type Row struct {
	keys.Key
	keys.GSI
	Shard

	// Type is the type of the row.
	UnmarshalledType string `dynamodbav:"type" json:"type,omitempty"`
	DBManager        `dynamodbav:"-" json:"-"`
}

// Type returns the type of the record.
func (r *Row) Type() string {
	if r.UnmarshalledType == "" {
		return "Row"
	}
	return r.UnmarshalledType
}

func (r *Row) Keys(gsi int) (string, string, error) {
	if r.PartitionKey != "" && r.SortKey != "" {
		return r.PartitionKey, r.SortKey, nil
	}
	if r.PartitionKey != "" && r.SortKey == "" {
		return r.PartitionKey, "row", nil
	}
	r.PartitionKey = uuid.UUIDv4()
	r.SortKey = "row"
	return r.PartitionKey, r.SortKey, nil
}
