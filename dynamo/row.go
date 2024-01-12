package dynamo

import (
	"encoding/json"
	"time"

	"github.com/dgryski/trifles/uuid"
)

// UnixTime represents a Unix timestamp in seconds.
type UnixTime time.Time

// MarshalDynamoDB converts UnixTime to a string representation for DynamoDB.
func (t UnixTime) MarshalDynamoDB() (string, error) {
	return time.Time(t).UTC().Format(time.RFC3339), nil
}

// UnmarshalDynamoDB converts a string representation from DynamoDB to UnixTime.
func (t *UnixTime) UnmarshalDynamoDB(s string) error {
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*t = UnixTime(parsedTime)
	return nil
}

// MarshalJSON converts UnixTime to a JSON representation.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

// UnmarshalJSON converts a JSON representation to UnixTime.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var unixTime int64
	err := json.Unmarshal(data, &unixTime)
	if err != nil {
		return err
	}
	*t = UnixTime(time.Unix(unixTime, 0))
	return nil
}

// AddTTL adds a duration to the UnixTime value. Negative durations are
// allowed, and will subtract from the UnixTime value.
func (t UnixTime) AddTTL(duration time.Duration) UnixTime {
	return UnixTime(time.Time(t).Add(duration))
}

// UpdateTTL updates the UnixTime value to the given time.
func (t *UnixTime) UpdateTTL(newTime time.Time) {
	*t = UnixTime(newTime)
}

// Row is a sample Keyable implementation. It is not intended to be used
// by itself, but rather to be embedded into other types. After embedding,
// you should implement the TableName and Keys methods on the parent type.
type Row struct {
	// Type string `dynamodbav:"type,omitempty" json:"type,omitempty"`
	Pk  string `dynamodbav:"pk,omitempty" json:"pk,omitempty"`
	Sk  string `dynamodbav:"sk,omitempty" json:"sk,omitempty"`
	Pk1 string `dynamodbav:"pk1,omitempty" json:"pk1,omitempty"`
	Sk1 string `dynamodbav:"sk1,omitempty" json:"sk1,omitempty"`
	Pk2 string `dynamodbav:"pk2,omitempty" json:"pk2,omitempty"`
	Sk2 string `dynamodbav:"sk2,omitempty" json:"sk2,omitempty"`
	Pk3 string `dynamodbav:"pk3,omitempty" json:"pk3,omitempty"`
	Sk3 string `dynamodbav:"sk3,omitempty" json:"sk3,omitempty"`
	Pk4 string `dynamodbav:"pk4,omitempty" json:"pk4,omitempty"`
	Sk4 string `dynamodbav:"sk4,omitempty" json:"sk4,omitempty"`
	Pk5 string `dynamodbav:"pk5,omitempty" json:"pk5,omitempty"`
	Sk5 string `dynamodbav:"sk5,omitempty" json:"sk5,omitempty"`
	Pk6 string `dynamodbav:"pk6,omitempty" json:"pk6,omitempty"`

	// TTL is the UTC time that this record will expire.
	TTL UnixTime `dynamodbav:"ttl,omitempty" json:"ttl,omitempty"`

	// PkShard is a field that is used
	PkShard string `dynamodbav:"pkshard,omitempty" json:"pkshard,omitempty"`
	// Type is the type of the row.
	UnmarshalledType string `dynamodbav:"type" json:"type,omitempty"`
	DBManager        `dynamodbav:"-" json:"-"`
}

// Type returns the type of the record.
func (r *Row) Type() string {
	if r.UnmarshalledType == "" {
		return "row"
	}
	return r.UnmarshalledType
}

// IsType returns true if the record is of the given type.
func (r Row) IsType(t string) bool {
	return r.Type() == t
}

func (r *Row) Keys(gsi int) (string, string, error) {
	if r.Pk != "" && r.Sk != "" {
		return r.Pk, r.Sk, nil
	}
	if r.Pk != "" && r.Sk == "" {
		return r.Pk, "row", nil
	}
	r.Pk = uuid.UUIDv4()
	r.Sk = "row"
	return r.Pk, r.Sk, nil
}

func (r *Row) MaxShard() int {
	return 100
}
