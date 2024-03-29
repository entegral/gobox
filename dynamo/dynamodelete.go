package dynamo

import (
	"context"

	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DeleteItem deletes a row from DynamoDB. The row must implement the Keyable
// interface. This method uses the default client. If you need to use a specific
// client, use DeleteItemWithClient instead, or use the client.SetDefaultClient method.
func (d *DBManager) DeleteItem(ctx context.Context, row types.Linkable) (*dynamodb.DeleteItemOutput, error) {
	if d.Client != nil {
		return d.deleteItemPrependTypeWithClient(ctx, d.Client, row)
	}
	return d.deleteItemPrependTypeWithClient(ctx, clients.GetDefaultClient(ctx), row)
}

func (d *DBManager) deleteItemPrependTypeWithClient(ctx context.Context, client *clients.Client, row types.Linkable) (*dynamodb.DeleteItemOutput, error) {
	pk, sk, err := row.Keys(0)
	if err != nil {
		return nil, err
	}
	properPk, err := addKeySegment(rowType, row.Type())
	if err != nil {
		return nil, err
	}
	seg, err := addKeySegment(rowPk, pk)
	if err != nil {
		return nil, err
	}
	properPk += seg
	rcc := awstypes.ReturnConsumedCapacityNone
	if checkTesting() {
		rcc = awstypes.ReturnConsumedCapacityTotal
	}
	key := map[string]awstypes.AttributeValue{
		"pk": &awstypes.AttributeValueMemberS{Value: properPk},
		"sk": &awstypes.AttributeValueMemberS{Value: sk},
	}
	tn := d.TableName(ctx)
	return client.Dynamo().DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:              &tn,
		Key:                    key,
		ReturnValues:           awstypes.ReturnValueAllOld,
		ReturnConsumedCapacity: rcc,
	})
}
