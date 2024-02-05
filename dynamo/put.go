package dynamo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Put[T Rowable](ctx context.Context, item T) (oldRow row[T], err error) {
	r := newRow[T](item)
	// Marshal the input into a map of AttributeValues
	rowData, err := attributevalue.MarshalMap(r.object)
	if err != nil {
		return oldRow, err
	}

	for i := 0; i < r.MaxGSIs(); i++ {
		pk, sk, err := r.object.Keys(i)
		if err != nil {
			msg := fmt.Sprintf("error generating keys for gsi %d of type %s", i, r.object.Type())
			slog.Error(msg, err)
			return oldRow, err
		}
		if pk == "" && sk == "" {
			continue
		} else if pk == "" {
			return oldRow, fmt.Errorf("partition key is required for gsi %d of type %s", i, r.object.Type())
		} else if sk == "" {
			return oldRow, fmt.Errorf("sort key is required for gsi %d of type %s", i, r.object.Type())
		}
		pkKey := "pk"
		skKey := "sk"
		if i > 0 {
			pkKey += fmt.Sprintf("%d", i)
			skKey += fmt.Sprintf("%d", i)
		}
		rowData[pkKey] = &awstypes.AttributeValueMemberS{Value: pk}
		rowData[skKey] = &awstypes.AttributeValueMemberS{Value: sk}
	}

	// Create the PutItem input
	putItemInput := &dynamodb.PutItemInput{
		Item:         rowData,
		TableName:    aws.String(r.TableName()),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB PutItem
	result, err := r.GetClient(ctx).Dynamo().PutItem(ctx, putItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}
