package dynamo

import (
	"context"

	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m *DiLink[T0, T1]) LoadEntity0s(ctx context.Context, linkWrapper ttypes.Typeable) ([]T0, error) {
	loaded, err := m.LoadEntity1(ctx)
	if err != nil {
		return nil, err
	}
	if !loaded {
		return nil, ErrEntityNotFound[T1]{Entity: m.Entity1}
	}

	links, err := FindByEntity1[T1, *DiLink[T0, T1]](ctx, m.Entity1, linkWrapper)
	if err != nil {
		return nil, err
	}
	var entities []T0
	for _, link := range links {
		loaded, err := link.LoadEntity0(ctx)
		if err != nil {
			return nil, err
		}
		if loaded {
			entities = append(entities, link.Entity0)
		}
	}
	return entities, nil
}

// findLinkRowsByEntity0 is a generic method to query for a list of rows based on the Entity0.
func findLinkRowsByEntity1[T1 ttypes.Linkable](ctx context.Context, clients *clients.Client, e1 T1, linkWrapper ttypes.Typeable) ([]map[string]types.AttributeValue, error) {
	return findLinkRowsByEntityGSI[T1](ctx, clients, e1, Entity1GSI, linkWrapper)
}

// FindByEntity1 is a generic method to query for a list of links based on the Entity1.
func FindByEntity1[T1, CustomLinkType ttypes.Linkable](ctx context.Context, e1 T1, linkWrapper ttypes.Typeable) ([]CustomLinkType, error) {
	clients := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntity1[T1](ctx, clients, e1, linkWrapper)
	if err != nil {
		return nil, err
	}
	var links []CustomLinkType
	for _, item := range rows {
		var link CustomLinkType
		if err := attributevalue.UnmarshalMap(item, &link); err != nil {
			return nil, err
		}
		// if err := validateDynamoRowType[CustomLinkType](item, link); err == nil {
		links = append(links, link)
		// }
	}
	return links, nil
}
