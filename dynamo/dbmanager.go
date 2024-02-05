package dynamo

import (
	"context"
	"encoding/json"
	"os"

	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// DynamoManager is a struct that implements the DynamoManager interface.
// It is intended to be embedded into other types to provide them with DynamoDB operations.
type DynamoManager[T types.Linkable] struct {
	Item  T
	Items []T
	// Tablename is a field you can set at runtime that will change the
	// name of the DynamoDB table that is used for operations.
	// It is ephemeral and will not be persisted to DynamoDB.
	tablename        string
	GetItemOutput    *dynamodb.GetItemOutput    `dynamodbav:"-" json:"-"`
	PutItemOutput    *dynamodb.PutItemOutput    `dynamodbav:"-" json:"-"`
	DeleteItemOutput *dynamodb.DeleteItemOutput `dynamodbav:"-" json:"-"`

	// RowData is a map of data retrieved from DynamoDB during the last
	// GetItem operation. This is useful for comparing the old values
	// with the new values after a PutItem operation.
	RowData map[string]awstypes.AttributeValue `dynamodbav:"-" json:"-"`
}

// DynamoManagerInterface defines the interface for DynamoDB operations.
// An interface is used to allow for mocking in unit tests, as well as to
// limit the scope of the methods that are exposed to the parent type.
type DynamoManagerInterface interface {
	TableName(ctx context.Context) string
	Get(ctx context.Context, row types.Linkable) (loaded bool, err error)
	WasGetSuccessful() bool
	Put(ctx context.Context, row types.Linkable) (err error)
	OldPutValues() map[string]awstypes.AttributeValue
	Delete(ctx context.Context, row types.Linkable) (err error)
	OldDeleteValues() map[string]awstypes.AttributeValue
	LoadFromMessage(ctx context.Context, message sqstypes.Message, row types.Linkable) (bool, error)
}

// NewDynamoManager creates a new instance of DynamoManager and returns it as a DynamoManagerInterface.
func NewDynamoManager(tableName string) *DynamoManager[types.Linkable] {
	return &DynamoManager[types.Linkable]{
		tablename: tableName,
	}
}

// TableName returns the name of the DynamoDB table.
// By default, this is the value of the TABLENAME environment variable.
// If you need to override this, implement this method on the parent type.
func (d *DynamoManager[T]) TableName(ctx context.Context) string {
	if d.tablename != "" {
		return d.tablename
	}
	tn := os.Getenv("TABLENAME")
	if tn == "" {
		panic("TABLENAME environment variable not set")
	}
	return tn
}

// SetTableName sets the name of the DynamoDB table.
func (d *DynamoManager[T]) SetTableName(tableName string) {
	d.tablename = tableName
}

// Get gets a row from DynamoDB. The row must implement the Keyable interface.
// The GetItemOutput response will be stored in the GetItemOutput field:
// d.GetItemOutput
func (d *DynamoManager[T]) Get(ctx context.Context) (item *T, err error) {
	tn := d.TableName(ctx)
	d.GetItemOutput, err = GetItemWithTablename(ctx, tn, d.Item)
	err = attributevalue.UnmarshalMap(d.GetItemOutput.Item, d.Item)
	if err != nil {
		return nil, err
	}
	return &d.Item, err
}

// Put puts an item into DynamoDB. The item must implement the Linkable interface.
// The PutItemOutput response will be stored in the PutItemOutput field:
// d.PutItemOutput
func (d *DynamoManager[T]) Put(ctx context.Context, item types.Linkable) (err error) {
	d.PutItemOutput, err = PutItem(ctx, item)
	return err
}

// OldPutValues returns the old values from the last successful PutItem operation.
func (d *DynamoManager[T]) OldPutValues() map[string]awstypes.AttributeValue {
	if d.PutItemOutput == nil {
		return nil
	}
	newItem := *d.PutItemOutput
	return newItem.Attributes
}

// Delete deletes a row from DynamoDB. The row must implement the Keyable interface.
// The DeleteItemOutput response will be stored in the DeleteItemOutput field:
// d.DeleteItemOutput
func (d *DynamoManager[T]) Delete(ctx context.Context) (err error) {
	d.DeleteItemOutput, err = DeleteItem(ctx, d.Item)
	return err
}

// OldDeleteValues returns the old values from the last successful DeleteItem operation.
func (d *DynamoManager[T]) OldDeleteValues() map[string]awstypes.AttributeValue {
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
func (d *DynamoManager[T]) LoadFromMessage(ctx context.Context, message sqstypes.Message) (*T, error) {
	if message.Body == nil || *message.Body == "" {
		return nil, &ErrSQSMessageEmpty{Message: message}
	}
	// Unmarshal the message body into the provided Row type
	if err := json.Unmarshal([]byte(*message.Body), d.Item); err != nil {
		return nil, err
	}

	// Use the existing Get method to load the item from DynamoDB
	return d.Get(ctx)
}
