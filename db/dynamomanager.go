package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/types"
)

type DynamoDBManager[T types.Keyable] struct {
	client    *dynamodb.Client
	tableName string
}

func (d *DynamoDBManager[T]) Get(ctx context.Context, input types.Keyable) (row Row[T], err error) {
	// Generate the keys from the input
	partitionKey, sortKey, err := input.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return row, err
	}

	// Create the key for DynamoDB
	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{
			Value: partitionKey,
		},
		"sk": &awstypes.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Create the GetItem input
	getItemInput := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(d.tableName),
	}

	// Call DynamoDB GetItem
	result, err := d.client.GetItem(ctx, getItemInput)
	if err != nil {
		return row, err
	}

	// Unmarshal the result into a Row
	err = attributevalue.UnmarshalMap(result.Item, &row)
	return row, err
}

func (d *DynamoDBManager[T]) Put(ctx context.Context, input types.Keyable) (oldRow Row[T], err error) {
	// Generate the keys from the input
	partitionKey, sortKey, err := input.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return oldRow, err
	}

	// Create the key for DynamoDB
	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{
			Value: partitionKey,
		},
		"sk": &awstypes.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Marshal the input into a map of AttributeValues
	item, err := attributevalue.MarshalMap(input)
	if err != nil {
		return oldRow, err
	}
	item["pk"] = &awstypes.AttributeValueMemberS{
		Value: partitionKey,
	}
	item["sk"] = &awstypes.AttributeValueMemberS{
		Value: sortKey,
	}
	// Create the PutItem input
	putItemInput := &dynamodb.PutItemInput{
		Item:         item,
		TableName:    aws.String(d.tableName),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB PutItem
	result, err := d.client.PutItem(ctx, putItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}

func (d *DynamoDBManager[T]) Delete(ctx context.Context, input types.Keyable) (oldRow Row[T], err error) {
	// Generate the keys from the input
	partitionKey, sortKey, err := input.Keys(0) // Assuming GSI 0 for primary key
	if err != nil {
		return oldRow, err
	}

	// Create the key for DynamoDB
	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{
			Value: partitionKey,
		},
		"sk": &awstypes.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Create the DeleteItem input
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key:          key,
		TableName:    aws.String(d.tableName),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB DeleteItem
	result, err := d.client.DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}
