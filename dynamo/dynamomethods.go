package dynamo

import (
	"context"
	"encoding/json"
	"os"

	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// DynamoDBOperations is a struct that implements the DynamoDBOperations interface.
// It is intended to be embedded into other types to provide them with DynamoDB operations.
type dynamoDBOperations struct {
	Tablename        string
	GetItemOutput    *dynamodb.GetItemOutput
	PutItemOutput    *dynamodb.PutItemOutput
	DeleteItemOutput *dynamodb.DeleteItemOutput

	// RowData is a map of the row data. This is used to store the raw data
	// from the dynamo response in the event the row data is needed for
	// custom, or subsequent, unmarshalling.
	RowData map[string]awstypes.AttributeValue `dynamodbav:"-" json:"-"`
}

// DynamoDBOperationsInterface defines the interface for DynamoDB operations.
type DynamoDBOperationsInterface interface {
	TableName(ctx context.Context) string
	Get(ctx context.Context, row types.Linkable) (loaded bool, err error)
	WasGetSuccessful() bool
	Put(ctx context.Context, row types.Linkable) (err error)
	OldPutValues(item any) map[string]awstypes.AttributeValue
	Delete(ctx context.Context, row types.Linkable) (err error)
	OldDeleteValues() map[string]awstypes.AttributeValue
	LoadFromMessage(ctx context.Context, message sqstypes.Message, row types.Linkable) (bool, error)
}

// NewDynamoDBOperations creates a new instance of DynamoDBOperations and returns it as a DynamoDBOperationsInterface.
func NewDynamoDBOperations(tableName string) DynamoDBOperationsInterface {
	return &dynamoDBOperations{
		Tablename: tableName,
	}
}

// TableName returns the name of the DynamoDB table.
// By default, this is the value of the TABLENAME environment variable.
// If you need to override this, implement this method on the parent type.
func (r *dynamoDBOperations) TableName(ctx context.Context) string {
	if r.Tablename != "" {
		return r.Tablename
	}
	tn := os.Getenv("TABLENAME")
	if tn == "" {
		panic("TABLENAME environment variable not set")
	}
	return tn
}

// Get gets a row from DynamoDB. The row must implement the Keyable interface.
// The GetItemOutput response will be stored in the GetItemOutput field:
// d.GetItemOutput
func (d *dynamoDBOperations) Get(ctx context.Context, row types.Linkable) (loaded bool, err error) {
	tn := d.TableName(ctx)
	d.GetItemOutput, err = GetItemWithTablename(ctx, tn, row)
	return d.WasGetSuccessful(), err
}

// WasGetSuccessful returns true if the last GetItem operation was successful.
func (d *dynamoDBOperations) WasGetSuccessful() bool {
	return d.GetItemOutput != nil && d.GetItemOutput.Item != nil
}

// Put puts a row into DynamoDB. The row must implement the Linkable interface.
// The PutItemOutput response will be stored in the PutItemOutput field:
// d.PutItemOutput
func (d *dynamoDBOperations) Put(ctx context.Context, row types.Linkable) (err error) {
	tn := d.TableName(ctx)
	d.PutItemOutput, err = PutItemWithTablename(ctx, tn, row)
	return err
}

// OldPutValues returns the old values from the last successful PutItem operation.
func (d *dynamoDBOperations) OldPutValues(item any) map[string]awstypes.AttributeValue {
	if d.PutItemOutput == nil {
		return nil
	}
	return d.PutItemOutput.Attributes
}

// Delete deletes a row from DynamoDB. The row must implement the Keyable interface.
// The DeleteItemOutput response will be stored in the DeleteItemOutput field:
// d.DeleteItemOutput
func (d *dynamoDBOperations) Delete(ctx context.Context, row types.Linkable) (err error) {
	d.DeleteItemOutput, err = DeleteItem(ctx, row)
	return err
}

// OldDeleteValues returns the old values from the last successful DeleteItem operation.
func (d *dynamoDBOperations) OldDeleteValues() map[string]awstypes.AttributeValue {
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
func (d *dynamoDBOperations) LoadFromMessage(ctx context.Context, message sqstypes.Message, row types.Linkable) (bool, error) {
	if message.Body == nil || *message.Body == "" {
		return false, ErrSQSMessageEmpty{Message: message}
	}
	// Unmarshal the message body into the provided Row type
	if err := json.Unmarshal([]byte(*message.Body), row); err != nil {
		return false, err
	}

	// Use the existing Get method to load the item from DynamoDB
	return d.Get(ctx, row)
}
