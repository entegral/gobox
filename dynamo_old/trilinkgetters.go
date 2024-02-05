package dynamo

import (
	"context"
	"reflect"

	"github.com/entegral/gobox/clients"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (m *TriLink[T0, T1, T2]) LoadEntity2(ctx context.Context) (bool, error) {
	var pk, sk, e2pk, e2sk string
	var err error

	if !reflect.ValueOf(m.Entity2).IsNil() {
		pk, sk, err = m.Entity2.Keys(0)
		if err != nil {
			return false, err
		}
		seg, err := addKeySegment(rowType, m.Entity2.Type())
		if err != nil {
			return false, err
		}
		e2pk = seg
		seg, err = addKeySegment(rowPk, pk)
		if err != nil {
			return false, err
		}
		e2pk += seg
		e2sk = sk
	} else {
		e2pk, e2sk, err = m.ExtractE2Keys()
		if err != nil {
			return false, err
		}
	}

	tn := m.TableName(ctx)
	clients := clients.GetDefaultClient(ctx)
	out, err := clients.Dynamo().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tn,
		Key: map[string]awstypes.AttributeValue{
			"pk": &awstypes.AttributeValueMemberS{Value: e2pk},
			"sk": &awstypes.AttributeValueMemberS{Value: e2sk},
		},
	})
	if err != nil {
		return false, err
	}
	if out.Item == nil {
		return false, &ErrEntityNotFound[T2]{Entity: m.Entity2}
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
