package dynamo

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dgryski/trifles/uuid"
)

func (r *Row) SetTTL(t time.Time) *UnixTime {
	r.TTL = &UnixTime{t}
	return r.TTL
}

// UnixTime represents a Unix timestamp in seconds.
type UnixTime struct {
	time.Time
}

// MarshalDynamoDBAttributeValue implements the dynamodbattribute.Marshaler interface.
func (t UnixTime) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return &types.AttributeValueMemberN{
		Value: strconv.FormatInt(t.UTC().Unix(), 10),
	}, nil
}

// UnmarshalDynamoDBAttributeValue implements the dynamodbattribute.Unmarshaler interface.
func (t *UnixTime) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	if v, ok := av.(*types.AttributeValueMemberN); ok {
		unixTime, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return err
		}
		t.Time = time.Unix(unixTime, 0)
	}
	return nil
}

// MarshalJSON converts UnixTime to a JSON representation.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.UTC().Format(time.RFC3339))
}

// UnmarshalJSON converts a JSON representation to UnixTime.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp string
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

// AddTTL adds a duration to the UnixTime value. Negative durations are
// allowed, and will subtract from the UnixTime value.
func (t UnixTime) Add(duration time.Duration) UnixTime {
	return UnixTime{t.Time.Add(duration)}
}

// UpdateTTL updates the UnixTime value to the given time.
func (t *UnixTime) UpdateTTL(newTime time.Time) {
	t.Time = newTime
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
	Sk6 string `dynamodbav:"sk6,omitempty" json:"sk6,omitempty"`

	// TTL is the UTC time that this record will expire.
	TTL *UnixTime `dynamodbav:"ttl,omitempty" json:"ttl,omitempty"`

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
