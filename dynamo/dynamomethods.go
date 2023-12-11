package dynamo

import (
	"context"

	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBOperations is a struct that implements the DynamoDBOperations interface.
// It is intended to be embedded into other types to provide them with DynamoDB operations.
type DynamoDBOperations struct {
	TableName        string
	GetItemOutput    *dynamodb.GetItemOutput
	PutItemOutput    *dynamodb.PutItemOutput
	DeleteItemOutput *dynamodb.DeleteItemOutput

	// RowData is a map of the row data. This is used to store the raw data
	// from the dynamo response in the event the row data is needed for
	// custom, or subsequent, unmarshalling.
	RowData map[string]awstypes.AttributeValue `dynamodbav:"-" json:"-"`
}

// Get gets a row from DynamoDB. The row must implement the Keyable interface.
// The GetItemOutput response will be stored in the GetItemOutput field:
// d.GetItemOutput
func (d *DynamoDBOperations) Get(ctx context.Context, row types.Linkable) (loaded bool, err error) {
	d.GetItemOutput, err = GetItem(ctx, row)
	return d.WasGetSuccessful(), err
}

// WasGetSuccessful returns true if the last GetItem operation was successful.
func (d *DynamoDBOperations) WasGetSuccessful() bool {
	return d.GetItemOutput != nil && d.GetItemOutput.Item != nil
}

// Put puts a row into DynamoDB. The row must implement the Linkable interface.
// The PutItemOutput response will be stored in the PutItemOutput field:
// d.PutItemOutput
func (d *DynamoDBOperations) Put(ctx context.Context, row types.Linkable) (err error) {
	d.PutItemOutput, err = PutItem(ctx, row)
	return err
}

// OldPutValues returns the old values from the last successful PutItem operation.
func (d *DynamoDBOperations) OldPutValues(item any) map[string]awstypes.AttributeValue {
	if d.PutItemOutput == nil {
		return nil
	}
	return d.PutItemOutput.Attributes
}

// func (d *DynamoDBOperations) Update(ctx context.Context, key map[string]awstypes.AttributeValue, updateExpression string, expressionAttributeValues map[string]awstypes.AttributeValue) (*dynamodb.UpdateItemOutput, (err error)) {
// 	return UpdateItem(ctx, d.DynamoDBClient, d.TableName, key, updateExpression, expressionAttributeValues)
// }

// Delete deletes a row from DynamoDB. The row must implement the Keyable interface.
// The DeleteItemOutput response will be stored in the DeleteItemOutput field:
// d.DeleteItemOutput
func (d *DynamoDBOperations) Delete(ctx context.Context, row types.Linkable) (err error) {
	d.DeleteItemOutput, err = DeleteItem(ctx, row)
	return err
}

// OldDeleteValues returns the old values from the last successful DeleteItem operation.
func (d *DynamoDBOperations) OldDeleteValues() map[string]awstypes.AttributeValue {
	if d.DeleteItemOutput == nil {
		return nil
	}
	return d.DeleteItemOutput.Attributes
}
