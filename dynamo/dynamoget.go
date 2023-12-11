package dynamo

import (
	"context"
	"reflect"

	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetItem gets a row from DynamoDB. The row must implement the Keyable
// interface. This method uses the default client. If you need to use a specific
// client, use GetItemWithClient instead, or use the client.SetDefaultClient method.
func GetItem(ctx context.Context, row types.Linkable) (*dynamodb.GetItemOutput, error) {
	return GetItemWithClient(ctx, clients.GetDefaultClient(ctx), row)
}

// GetItemWithClient gets a row from DynamoDB using the provided client
// The row must implement the Keyable interface
func GetItemWithClient(ctx context.Context, client *clients.Client, row types.Linkable) (*dynamodb.GetItemOutput, error) {
	pk, sk := row.Keys(0)

	pkWithTypePrefix := addKeySegment(rowType, row.Type())
	pkWithTypePrefix += addKeySegment(rowPk, pk)

	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{Value: pkWithTypePrefix},
		"sk": &awstypes.AttributeValueMemberS{Value: sk},
	}

	out, err := client.Dynamo().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(clients.TableName(ctx)),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	// var newRow T

	// Check if T implements CustomDynamoMarshaller
	if marshaller, ok := any(row).(types.CustomDynamoMarshaller); ok {
		err = marshaller.UnmarshalItem(out.Item)
	} else {
		err = attributevalue.UnmarshalMap(out.Item, row)
	}
	if err != nil {
		return nil, err
	}
	// if the row has a RowData field by embedding the Row struct, set it
	rowValue := reflect.ValueOf(row).Elem()
	if rowDataField := rowValue.FieldByName("Row"); rowDataField.IsValid() {
		if dynamoDBOperations := rowDataField.FieldByName("DynamoDBOperations"); dynamoDBOperations.IsValid() {
			rowData := dynamoDBOperations.FieldByName("RowData")
			if rowData.CanSet() {
				rowData.Set(reflect.ValueOf(out.Item))
			}
		}
	}
	return out, nil
}
