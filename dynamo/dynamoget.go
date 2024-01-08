package dynamo

import (
	"context"
	"os"
	"reflect"

	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetItem gets a row from DynamoDB. The row must implement the Keyable
// interface.
func GetItem(ctx context.Context, row types.Linkable) (*dynamodb.GetItemOutput, error) {
	tn := types.CheckTableable(ctx, row)
	return getItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), tn, row)
}

func GetItemWithTablename(ctx context.Context, tablename string, row types.Linkable) (*dynamodb.GetItemOutput, error) {
	return getItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), tablename, row)
}

func getItemPrependTypeWithClient(ctx context.Context, client *clients.Client, tablename string, row types.Linkable) (*dynamodb.GetItemOutput, error) {
	pk, sk, err := row.Keys(0)
	if err != nil {
		return nil, err
	}

	pkWithTypePrefix, err := addKeySegment(rowType, row.Type())
	if err != nil {
		return nil, err
	}
	seg, err := addKeySegment(rowPk, pk)
	if err != nil {
		return nil, err
	}
	pkWithTypePrefix += seg

	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{Value: pkWithTypePrefix},
		"sk": &awstypes.AttributeValueMemberS{Value: sk},
	}
	rcc := awstypes.ReturnConsumedCapacityNone
	if os.Getenv("TESTING") == "true" {
		rcc = awstypes.ReturnConsumedCapacityTotal
	}
	out, err := client.Dynamo().GetItem(ctx, &dynamodb.GetItemInput{
		TableName:              aws.String(tablename),
		Key:                    key,
		ReturnConsumedCapacity: rcc,
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
		if dynamoDBOperations := rowDataField.FieldByName("dynamoDBOperations"); dynamoDBOperations.IsValid() {
			rowData := dynamoDBOperations.FieldByName("RowData")
			if rowData.CanSet() {
				rowData.Set(reflect.ValueOf(out.Item))
			}
		}
	}
	return out, nil
}
