package dynamo

import (
	"context"

	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m *TriLink[T0, T1, T2]) LoadEntity2(ctx context.Context) (bool, error) {
	e1pk, e1sk := m.ExtractE2Keys()
	tn := m.TableName(ctx)
	clients := clients.GetDefaultClient(ctx)
	out, err := clients.Dynamo().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tn,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: e1pk},
			"sk": &types.AttributeValueMemberS{Value: e1sk},
		},
	})
	if err != nil {
		return false, err
	}
	if out.Item == nil {
		return false, ErrEntityNotFound[T2]{Entity: m.Entity2}
	}
	if err := validateDynamoRowType[T2](out.Item, m.Entity2); err != nil {
		return false, err
	}
	err = attributevalue.UnmarshalMap(out.Item, &m.Entity2)
	if err != nil {
		return false, err
	}
	return true, nil
}

// findLinkRowsByEntity0 is a generic method to query for a list of rows based on the Entity0.
func findLinkRowsByEntity2[T2 ttypes.Linkable](ctx context.Context, clients *clients.Client, e2 T2) ([]map[string]types.AttributeValue, error) {
	return findLinkRowsByEntityGSI[T2](ctx, clients, e2, Entity2GSI)
}

// FindByEntity2 is a generic method to query for a list of links based on the Entity2.
func FindByEntity2[T2, CustomLinkType ttypes.Linkable](ctx context.Context, e1 T2) ([]CustomLinkType, error) {
	clients := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntity2[T2](ctx, clients, e1)
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
