package dynamo

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"gobox/clients"
	"gobox/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Query executes a query against DynamoDB. The method uses the default client.
// If you need to use a specific client, use QueryStandardUnmarshal instead, or use the
// client.SetDefaultClient method to set the default client.
//
// It also uses the standard attributevalue.UnmarshalMap to unmarshal the results.
// If you need to use a custom unmarshaler, satisfy the interface requred for/used by the
// QueryCustomUnmarshal function.
func Query[T any](ctx context.Context, input dynamodb.QueryInput) ([]T, error) {
	return QueryListOfMaps[T](ctx, clients.GetDefaultClient(ctx), input)
}

// QueryByGSI executes a query against a specific Global Secondary Index (GSI) or the primary index.
// It uses the default client. If you need to use a specific client, use QueryByGSIWithClient instead.
func QueryByGSI[T types.Keyable](ctx context.Context, gsi int, partitionKey, sortKeyCondition string) ([]T, error) {
	client := clients.GetDefaultClient(ctx)
	return QueryByGSIWithClient[T](ctx, client, gsi, partitionKey, sortKeyCondition)
}

// QueryByGSIWithClient executes a query against a specific Global Secondary Index (GSI) or the primary index using the provided client.
func QueryByGSIWithClient[T types.Keyable](ctx context.Context, client *clients.Client, gsi int, partitionKey, sortKeyCondition string) ([]T, error) {
	tableName := os.Getenv("TABLENAME")
	if tableName == "" {
		return nil, fmt.Errorf("TABLENAME environment variable not set")
	}

	var indexName *string
	var keyConditionExpression string
	var expressionAttributeNames map[string]string
	var expressionAttributeValues map[string]awstypes.AttributeValue

	// Construct the key condition expression based on the GSI
	if gsi == 0 {
		keyConditionExpression = "pk = :pk"
		if sortKeyCondition != "" {
			keyConditionExpression += " AND begins_with(sk, :sk)"
		}
	} else {
		indexName = aws.String(fmt.Sprintf("pk%d-sk%d-index", gsi, gsi))
		keyConditionExpression = fmt.Sprintf("pk%d = :pk", gsi)
		if sortKeyCondition != "" {
			keyConditionExpression += fmt.Sprintf(" AND begins_with(sk%d, :sk)", gsi)
		}
	}

	// Construct the expression attribute values
	expressionAttributeValues = map[string]awstypes.AttributeValue{
		":pk": &awstypes.AttributeValueMemberS{Value: partitionKey},
	}
	if sortKeyCondition != "" {
		expressionAttributeValues[":sk"] = &awstypes.AttributeValueMemberS{Value: sortKeyCondition}
	}

	// Create the query input
	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 indexName,
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	// Execute the query
	return QueryListOfMaps[T](ctx, client, queryInput)
}

// QueryListOfMaps executes a query against DynamoDB and unmarshals using attributevalue.UnmarshalMap.
func QueryListOfMaps[T any](ctx context.Context, client *clients.Client, input dynamodb.QueryInput) ([]T, error) {
	out, err := client.Dynamo().Query(ctx, &input)
	if err != nil {
		return nil, err
	}

	var rows []T
	err = attributevalue.UnmarshalListOfMaps(out.Items, &rows)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		// Use reflection to check if row embeds the Row type and update RowData
		rowValue := reflect.ValueOf(row).Elem()
		if rowDataField := rowValue.FieldByName("Row"); rowDataField.IsValid() {
			// Set the RowData field with the raw data from DynamoDB
			rowData := rowDataField.FieldByName("RowData")
			if rowData.CanSet() {
				rowData.Set(reflect.ValueOf(out.Items[i]))
			}
		}
	}
	return rows, nil
}

// QueryCustomListOfMaps executes a query against DynamoDB and unmarshals using the custom interface.
func QueryCustomListOfMaps[T types.CustomDynamoMarshaller](ctx context.Context, client *clients.Client, input dynamodb.QueryInput) ([]T, error) {
	out, err := client.Dynamo().Query(ctx, &input)
	if err != nil {
		return nil, err
	}

	var rows []T
	for _, item := range out.Items {
		var newRow T
		err := newRow.UnmarshalItem(item)
		if err != nil {
			return nil, err
		}
		// Use reflection to check if row embeds the Row type and update RowData
		rowValue := reflect.ValueOf(newRow).Elem()
		if rowDataField := rowValue.FieldByName("Row"); rowDataField.IsValid() {
			// Set the RowData field with the raw data from DynamoDB
			rowData := rowDataField.FieldByName("RowData")
			if rowData.CanSet() {
				rowData.Set(reflect.ValueOf(item))
			}
		}
		rows = append(rows, newRow)
	}
	return rows, nil
}
