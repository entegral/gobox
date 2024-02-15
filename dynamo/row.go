package dynamo

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dgryski/trifles/uuid"
	"github.com/entegral/gobox/keys"
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
	keys.Key
	keys.GSI
	Shard

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
