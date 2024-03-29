package dynamo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// DBManager is a struct that implements the DBManager interface.
// It is intended to be embedded into other types to provide them with DynamoDB operations.
type DBManager struct {
	Table
	// Tablename is a field you can set at runtime that will change the
	// name of the DynamoDB table that is used for operations.
	// It is ephemeral and will not be persisted to DynamoDB.
	Tablename        string
	GetItemOutput    *dynamodb.GetItemOutput    `dynamodbav:"-" json:"-"`
	PutItemOutput    *dynamodb.PutItemOutput    `dynamodbav:"-" json:"-"`
	DeleteItemOutput *dynamodb.DeleteItemOutput `dynamodbav:"-" json:"-"`

	// TTL is a UnixTime timestamp that is used to set the Time To Live
	// (TTL) for the item in DynamoDB.
	TTL UnixTime `dynamodbav:"ttl,omitempty" json:"ttl,omitempty"`

	// RowData is a map of data retrieved from DynamoDB during the last
	// GetItem operation. This is useful for comparing the old values
	// with the new values after a PutItem operation.
	RowData map[string]awstypes.AttributeValue `dynamodbav:"-" json:"-"`
}

// DBManagerInterface defines the interface for DynamoDB operations.
// An interface is used to allow for mocking in unit tests, as well as to
// limit the scope of the methods that are exposed to the parent type.
type DBManagerInterface interface {
	TableName(ctx context.Context) string
	Get(ctx context.Context, row types.Linkable) (loaded bool, err error)
	WasGetSuccessful() bool
	Put(ctx context.Context, row types.Linkable) (err error)
	OldPutValues() map[string]awstypes.AttributeValue
	Delete(ctx context.Context, row types.Linkable) (err error)
	OldDeleteValues() map[string]awstypes.AttributeValue
	LoadFromMessage(ctx context.Context, message sqstypes.Message, row types.Linkable) (bool, error)
	SetDynamoTTL(t time.Time) *UnixTime
	GetDynamoTTL() *UnixTime
}

func (d *DBManager) SetDynamoTTL(t time.Time) *UnixTime {
	d.TTL = UnixTime{t}
	return &d.TTL
}

func (d *DBManager) GetDynamoTTL() *UnixTime {
	return &d.TTL
}

// NewDBManager creates a new instance of DBManager and returns it as a DBManagerInterface.
func NewDBManager(tableName string) DBManagerInterface {
	return &DBManager{
		Tablename: tableName,
	}
}

// Get gets a row from DynamoDB. The row must implement the Keyable interface.
// The GetItemOutput response will be stored in the GetItemOutput field:
// d.GetItemOutput
func (d *DBManager) Get(ctx context.Context, row types.Linkable) (loaded bool, err error) {
	d.GetItemOutput, err = d.GetItemWithTablename(ctx, row)
	return d.WasGetSuccessful(), err
}

// WasGetSuccessful returns true if the last GetItem operation was successful.
func (d *DBManager) WasGetSuccessful() bool {
	return d.GetItemOutput != nil && d.GetItemOutput.Item != nil
}

// Put puts a row into DynamoDB. The row must implement the Linkable interface.
// The PutItemOutput response will be stored in the PutItemOutput field:
// d.PutItemOutput
func (d *DBManager) Put(ctx context.Context, row types.Linkable) (err error) {
	d.PutItemOutput, err = d.PutItem(ctx, row)
	return err
}

func (d *DBManager) WasPutSuccessful() bool {
	return d.PutItemOutput != nil
}

// OldPutValues returns the old values from the last successful PutItem operation.
func (d *DBManager) OldPutValues() map[string]awstypes.AttributeValue {
	if d.PutItemOutput == nil {
		return nil
	}
	newItem := *d.PutItemOutput
	return newItem.Attributes
}

// Delete deletes a row from DynamoDB. The row must implement the Keyable interface.
// The DeleteItemOutput response will be stored in the DeleteItemOutput field:
// d.DeleteItemOutput
func (d *DBManager) Delete(ctx context.Context, row types.Linkable) (err error) {
	d.DeleteItemOutput, err = d.DeleteItem(ctx, row)
	return err
}

// OldDeleteValues returns the old values from the last successful DeleteItem operation.
func (d *DBManager) OldDeleteValues() map[string]awstypes.AttributeValue {
	if d.DeleteItemOutput == nil {
		return nil
	}
	return d.DeleteItemOutput.Attributes
}

type ErrSQSMessageEmpty struct {
	Message sqstypes.Message
}

func (e ErrSQSMessageEmpty) Error() string {
	return "sqs message body is empty"
}

// LoadFromMessage unmarshals an SQS message into a Row and then loads the full item from DynamoDB.
func (d *DBManager) LoadFromMessage(ctx context.Context, message sqstypes.Message, row types.Linkable) (bool, error) {
	if message.Body == nil || *message.Body == "" {
		return false, &ErrSQSMessageEmpty{Message: message}
	}
	// Unmarshal the message body into the provided Row type
	if err := json.Unmarshal([]byte(*message.Body), row); err != nil {
		return false, err
	}

	// Use the existing Get method to load the item from DynamoDB
	return d.Get(ctx, row)
}
