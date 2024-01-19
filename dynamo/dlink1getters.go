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
	links, err := FindLinksByEntity0[T0, *DiLink[T0, T1]](ctx, m.Entity0, linkWrapper)
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

func (m *DiLink[T0, T1]) LoadEntity1(ctx context.Context) (bool, error) {
	pk, sk, err := m.Entity1.Keys(0)
	var e1pk, e1sk string
	if err == nil {
		seg, err := addKeySegment(rowType, m.Entity1.Type())
		if err != nil {
			return false, err
		}
		e1pk = seg
		seg, err = addKeySegment(rowPk, pk)
		if err != nil {
			return false, err
		}
		e1pk += seg
		e1sk = sk
	} else {
		if pk == "" || sk == "" {
			e1pk, e1sk, err = m.ExtractE1Keys()
			if err != nil {
				return false, err
			}
		}
	}
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
