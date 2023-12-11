package dynamo

import (
	"context"

	"gobox/clients"
	"gobox/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// UpdateItem updates a row in DynamoDB. The row must implement the
// DynamoUpdater interface. Consider embedding your type into a wrapper
// that implements DynamoUpdater in order to issue the desired update behavior.
//
// This method uses the default client. If you need to use a specific client,
// use UpdateItemWithClient instead, or use the client.SetDefaultClient method.
func UpdateItem(ctx context.Context, row types.DynamoUpdater) (*dynamodb.UpdateItemOutput, error) {
	client := clients.GetDefaultClient(ctx)
	return UpdateItemWithClient(ctx, client, row)
}

// UpdateItemWithClient updates a row in DynamoDB. The row must implement the
// DynamoUpdater interface. Consider embedding your type into a wrapper
// that implements DynamoUpdater in order to issue the desired update behavior.
func UpdateItemWithClient(ctx context.Context, client *clients.Client, row types.DynamoUpdater) (*dynamodb.UpdateItemOutput, error) {
	input, err := row.DynamoUpdateInput(ctx)
	if err != nil {
		return nil, err
	}
	return client.Dynamo().UpdateItem(ctx, input)
}
