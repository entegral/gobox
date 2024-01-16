package dynamo

import (
	"context"

	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m *DiLink[T0, T1]) LoadEntity1s(ctx context.Context, linkWrapper ttypes.Typeable) ([]T1, error) {
	loaded, err := m.LoadEntity0(ctx)
	if err != nil {
		return nil, err
	}
	if !loaded {
		return nil, ErrEntityNotFound[T0]{Entity: m.Entity0}
	}
	links, err := FindByEntity0[T0, *DiLink[T0, T1]](ctx, m.Entity0, linkWrapper)
	if err != nil {
		return nil, err
	}
	var entities []T1
	for _, link := range links {
		loaded, err := link.LoadEntity1(ctx)
		if err != nil {
			return nil, err
		}
		if loaded {
			entities = append(entities, link.Entity1)
		}
	}
	return entities, nil
}

// FindByEntity0 is a generic method to query for a list of links based on the Entity0.
func FindByEntity0[T0, CustomLinkType ttypes.Linkable](ctx context.Context, e0 T0, linkWrapper ttypes.Typeable) ([]CustomLinkType, error) {
	client := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntityGSI[T0](ctx, client, e0, Entity0GSI, linkWrapper)
	if err != nil {
		return nil, err
	}
	var links []CustomLinkType
	for _, item := range rows {
		var link CustomLinkType
		if err := attributevalue.UnmarshalMap(item, &link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (m *DiLink[T0, T1]) LoadEntity1(ctx context.Context) (bool, error) {
	e1pk, e1sk := m.ExtractE1Keys()
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
		return false, &ErrEntityNotFound[T1]{Entity: m.Entity1}
	}
	if err := validateDynamoRowType[T1](out.Item, m.Entity1); err != nil {
		return false, err
	}
	err = attributevalue.UnmarshalMap(out.Item, &m.Entity1)
	if err != nil {
		return false, err
	}
	return true, nil
}
