package types

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dateTimeFormat = time.RFC3339

// SetDefaultDateTimeFormat sets the default format for DateTime values.
func SetDefaultDateTimeFormat(format string) {
	dateTimeFormat = format
}

// ResetDefaultDateTimeFormat resets the default format for DateTime values to RFC3339.
func ResetDefaultDateTimeFormat() {
	dateTimeFormat = time.RFC3339
}

// DateTime is a custom scalar type for handling date-time values.
type DateTime struct {
	time.Time
}

func (d DateTime) String() string {
	return d.Time.Format(dateTimeFormat)
}

// MarshalGQL implements the graphql.Marshaler interface.
func (d DateTime) MarshalGQL(w io.Writer) {
	w.Write([]byte(d.Time.Format(dateTimeFormat)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (d *DateTime) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("DateTime must be a string")
	}

	var err error
	d.Time, err = time.Parse(dateTimeFormat, str)
	return err
}

// MarshalJSON for JSON serialization.
func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(dateTimeFormat) + `"`), nil
}

// UnmarshalJSON for JSON deserialization.
func (d *DateTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"`+dateTimeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalDynamoDBAttributeValue for DynamoDB serialization.
func (d DateTime) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: d.Time.Format(dateTimeFormat)}, nil
}

// UnmarshalDynamoDBAttributeValue for DynamoDB deserialization.
func (d *DateTime) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	sValue, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("DateTime must be a string in DynamoDB")
	}

	var err error
	d.Time, err = time.Parse(dateTimeFormat, sValue.Value)
	if err != nil {
		return err
	}

	return nil
}
