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

// LoadEntity0 attempts to load the Entity0 from DynamoDB.
// If the Entity0 field is already populated, it will generate the keys
// and attempt to load the Entity0 from DynamoDB.
//
// If the Entity0 field is not populated, it will attempt to extract the keys
// from the Pk and Sk fields and then load the Entity0 from DynamoDB.
func (m *MonoLink[T0]) LoadEntity0(ctx context.Context) (bool, error) {
	pk, sk, err := m.Entity0.Keys(0)
	var e0pk, e0sk string
	if err == nil {
		seg, err := addKeySegment(rowType, m.Entity0.Type())
		if err != nil {
			return false, err
		}
		e0pk = seg
		seg, err = addKeySegment(rowPk, pk)
		if err != nil {
			return false, err
		}
		e0pk += seg
		e0sk = sk
	} else {
		if pk == "" || sk == "" {
			e0pk, e0sk, err = m.ExtractE0Keys()
			if err != nil {
				return false, err
			}
		}
	}
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
		return false, &ErrEntityNotFound[T0]{Entity: m.Entity0}
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
func findLinkRowsByEntity0[T0 ttypes.Linkable](ctx context.Context, clients *clients.Client, e0 T0, linkWrapper ttypes.Typeable) ([]map[string]types.AttributeValue, error) {
	return findLinkRowsByEntityGSI[T0](ctx, clients, e0, Entity0GSI, linkWrapper)
}

// findLinkRowsByEntityGSI is a generic method to query for a list of rows based on the Entity1.
func findLinkRowsByEntityGSI[T ttypes.Linkable](ctx context.Context, clients *clients.Client, entity T, entityGSI EntityGSI, linkWrapper ttypes.Typeable) ([]map[string]types.AttributeValue, error) {
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

	ePk, eSk, err := entity.Keys(0)
	if err != nil {
		return nil, err
	}
	linkedPk, err := addKeySegment(rowType, entity.Type())
	if err != nil {
		return nil, err
	}
	seg, err := addKeySegment(rowPk, ePk)
	if err != nil {
		return nil, err
	}
	linkedPk += seg

	kce := fmt.Sprintf("%s = :pk AND begins_with(%s, :sk)", epkKey, eskKey)
	tn := ttypes.CheckTableable(ctx, entity)
	index := entityGSI.String()
	fe := "#type = :type"
	qi := dynamodb.QueryInput{
		TableName:              &tn,
		KeyConditionExpression: &kce,
		IndexName:              &index,
		FilterExpression:       &fe,
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":   &types.AttributeValueMemberS{Value: linkedPk},
			":sk":   &types.AttributeValueMemberS{Value: eSk},
			":type": &types.AttributeValueMemberS{Value: linkWrapper.Type()},
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
