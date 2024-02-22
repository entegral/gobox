package dynamo

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// UnixTime represents a Unix timestamp in seconds.
type UnixTime struct {
	time.Time
}

// MarshalAttribute implements the attributevalue.Marshaler interface.
func (t UnixTime) MarshalAttribute() (types.AttributeValue, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return &types.AttributeValueMemberN{
		Value: strconv.FormatInt(t.UTC().Unix(), 10),
	}, nil
}

// UnmarshalAttribute implements the attributevalue.Unmarshaler interface.
func (t *UnixTime) UnmarshalAttribute(av types.AttributeValue) error {
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
