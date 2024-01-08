package dynamo

import (
	"context"
	"log/slog"

	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
		return false, ErrEntityNotFound[T1]{Entity: m.Entity1}
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

func (m *DiLink[T0, T1]) LoadEntity1s(ctx context.Context, linkWrapper ttypes.Typeable) ([]T1, error) {
	client := clients.GetDefaultClient(ctx)
	loaded, err := m.LoadEntity0(ctx)
	if err != nil {
		return nil, err
	}
	if !loaded {
		return nil, ErrEntityNotFound[T0]{Entity: m.Entity0}
	}
	return findEntity1s[T0, T1](ctx, client, m.Entity0, linkWrapper)
}

// findLinksByEntity0 is a generic method to query for a list of links based on the Entity0.
func findDiLinksByEntity0[T0, T1 ttypes.Linkable](ctx context.Context, clients *clients.Client, e0 T0, linkWrapper ttypes.Typeable) ([]DiLink[T0, T1], error) {
	rows, err := findLinkRowsByEntityGSI[T0](ctx, clients, e0, Entity0GSI, linkWrapper)
	if err != nil {
		return nil, err
	}
	var links []DiLink[T0, T1]
	err = attributevalue.UnmarshalListOfMaps(rows, &links)
	if err != nil {
		return nil, err
	}
	return links, nil
}

// findEntity1s is a generic method to load a list of Entity1s linked to the provided Entity0.
func findEntity1s[T0, T1 ttypes.Linkable](ctx context.Context, clients *clients.Client, e0 T0, linkWrapper ttypes.Typeable) ([]T1, error) {
	links, err := findDiLinksByEntity0[T0, T1](ctx, clients, e0, linkWrapper)
	if err != nil {
		return nil, err
	}
	var entity1s []T1
	for _, link := range links {
		link.Entity0 = e0
		loaded, err := link.LoadEntity1(ctx)
		if err != nil {
			slog.Error("error loading entity1 for link", err)
			continue
		}
		if loaded {
			entity1s = append(entity1s, link.Entity1)
		}
	}
	return entity1s, nil
}
