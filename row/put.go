package row

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Modify the Put function to use the GenerateKeys function
func (item *Row[T]) Put(ctx context.Context) (oldRow Row[T], err error) {
	// Marshal the input into a map of AttributeValues
	rowData, err := attributevalue.MarshalMap(item.object)
	if err != nil {
		return oldRow, err
	}

	// Create channels for the keys, post-processed keys, and errors
	keys := make(chan Key)
	processedKeys := make(chan Key)
	errs := make(chan error)

	// Start a goroutine to generate the keys
	go item.GenerateKeys(ctx, keys, errs)

	// Start a goroutine to post-process the keys
	go item.postProcessKeys(ctx, keys, processedKeys, errs)

	// Process the post-processed keys and errors
	for key := range processedKeys {
		pkKey := "pk"
		skKey := "sk"
		if key.Index > 0 {
			pkKey += fmt.Sprintf("%d", key.Index)
			skKey += fmt.Sprintf("%d", key.Index)
		}
		rowData[pkKey] = &awstypes.AttributeValueMemberS{Value: key.PK}
		rowData[skKey] = &awstypes.AttributeValueMemberS{Value: key.SK}

		// ensure the item's key struct is updated with the new keys
		item.Keys.SetKey(key)
	}

	// Check for any errors
	for err := range errs {
		return oldRow, err
	}

	// Create the PutItem input
	putItemInput := &dynamodb.PutItemInput{
		Item:         rowData,
		TableName:    aws.String(item.TableName()),
		ReturnValues: awstypes.ReturnValueAllOld,
	}

	// Call DynamoDB PutItem
	result, err := item.GetClient(ctx).Dynamo().PutItem(ctx, putItemInput)
	if err != nil {
		return oldRow, err
	}

	// Unmarshal the old item into a Row
	err = attributevalue.UnmarshalMap(result.Attributes, &oldRow)
	return oldRow, err
}
