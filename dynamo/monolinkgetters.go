package dynamo

import (
	"context"
	"fmt"

	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m *MonoLink[T0]) LoadEntity0(ctx context.Context) (bool, error) {
	e0pk, e0sk := m.ExtractE0Keys()
	tn := m.TableName(ctx)
	clients := clients.GetDefaultClient(ctx)
	out, err := clients.Dynamo().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tn,
		Key: map[string]awstypes.AttributeValue{
			"pk": &awstypes.AttributeValueMemberS{Value: e0pk},
			"sk": &awstypes.AttributeValueMemberS{Value: e0sk},
		},
	})
	if err != nil {
		return false, err
	}
	if out.Item == nil {
		return false, ErrEntityNotFound[T0]{Entity: m.Entity0}
	}
	if err := validateDynamoRowType[T0](out.Item, m.Entity0); err != nil {
		return false, err
	}
	err = attributevalue.UnmarshalMap(out.Item, &m.Entity0)
	if err != nil {
		return false, err
	}
	return true, nil
}

// findLinkRowsByEntity0 is a generic method to query for a list of rows based on the Entity0.
func findLinkRowsByEntity0[T0 ttypes.Linkable](ctx context.Context, e0 T0) ([]map[string]types.AttributeValue, error) {
	client := clients.GetDefaultClient(ctx)
	return findLinkRowsByEntityGSI[T0](ctx, client, e0, Entity0GSI)
}

// FindCustomLinksByEntity1 is a generic method to query for a list of links based on the Entity1.
func FindByEntity0[T0, CustomLinkType ttypes.Linkable](ctx context.Context, e0 T0) ([]CustomLinkType, error) {
	rows, err := findLinkRowsByEntity0[T0](ctx, e0)
	if err != nil {
		return nil, err
	}
	var links []CustomLinkType
	for _, item := range rows {
		var link CustomLinkType
		if err := attributevalue.UnmarshalMap(item, &link); err != nil {
			return nil, err
		}
		if err := validateDynamoRowType[CustomLinkType](item, link); err == nil {
			links = append(links, link)
		}
	}
	return links, nil
}

// findLinkRowsByEntityGSI is a generic method to query for a list of rows based on the Entity1.
func findLinkRowsByEntityGSI[T ttypes.Linkable](ctx context.Context, clients *clients.Client, entity T, entityGSI EntityGSI) ([]map[string]types.AttributeValue, error) {
	var epkKey, eskKey string
	switch entityGSI {
	case Entity0GSI:
		epkKey = entity0pk.String()
		eskKey = entity0sk.String()
	case Entity1GSI:
		epkKey = entity1pk.String()
		eskKey = entity1sk.String()
	case Entity2GSI:
		epkKey = entity2pk.String()
		eskKey = entity2sk.String()
	}

	ePk, eSk := entity.Keys(0)
	kce := fmt.Sprintf("%s = :pk AND begins_with(%s, :sk)", epkKey, eskKey)
	tn := clients.TableName(ctx)
	index := entityGSI.String()
	qi := dynamodb.QueryInput{
		TableName:              &tn,
		KeyConditionExpression: &kce,
		IndexName:              &index,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: ePk},
			":sk": &types.AttributeValueMemberS{Value: eSk},
		},
	}
	out, err := clients.Dynamo().Query(ctx, &qi)
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, nil
	}
	return out.Items, nil
}

// Get gets a row from DynamoDB. The row must implement the Keyable interface.
// The GetItemOutput response will be stored in the GetItemOutput field:
// d.GetItemOutput
func (m *MonoLink[T0]) Get(ctx context.Context, row ttypes.Linkable) (err error) {
	m.GetItemOutput, err = GetItemPrependType(ctx, row)
	return err
}

// WasGetSuccessful returns true if the last GetItem operation was successful.
func (m *MonoLink[T0]) WasGetSuccessful() bool {
	return m.GetItemOutput != nil && m.GetItemOutput.Item != nil
}

// Put puts a row into DynamoDB. The row must implement the Linkable interface.
// The PutItemOutput response will be stored in the PutItemOutput field:
// m.PutItemOutput
func (m *MonoLink[T0]) Put(ctx context.Context, row ttypes.Linkable) (err error) {
	m.PutItemOutput, err = PutItemPrependType(ctx, row)
	return err
}

// OldPutValues returns the old values from the last successful PutItem operation.
func (m *MonoLink[T0]) OldPutValues(item any) map[string]awstypes.AttributeValue {
	if m.PutItemOutput == nil {
		return nil
	}
	return m.PutItemOutput.Attributes
}

// func (m *MonoLink[T0]) Update(ctx context.Context, key map[string]awstypes.AttributeValue, updateExpression string, expressionAttributeValues map[string]awstypes.AttributeValue) (*dynamodb.UpdateItemOutput, (err error)) {
// 	return UpdateItem(ctx, m.DynamoDBClient, m.TableName, key, updateExpression, expressionAttributeValues)
// }

// Delete deletes a row from DynamoDB. The row must implement the Keyable interface.
// The DeleteItemOutput response will be stored in the DeleteItemOutput field:
// m.DeleteItemOutput
func (m *MonoLink[T0]) Delete(ctx context.Context, row ttypes.Linkable) (err error) {
	m.DeleteItemOutput, err = DeleteItemPrependType(ctx, row)
	return err
}
