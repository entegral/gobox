package dynamo

import (
	"context"
	"reflect"

	"github.com/entegral/gobox/clients"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// LoadEntity0 attempts to load the Entity0 from DynamoDB.
// If the Entity0 field is already populated, it will generate the keys
// and attempt to load the Entity0 from DynamoDB.
//
// If the Entity0 field is not populated, it will attempt to extract the keys
// from the Pk and Sk fields and then load the Entity0 from DynamoDB.
func (m *MonoLink[T0]) LoadEntity0(ctx context.Context) (bool, error) {
	var pk, sk, e0pk, e0sk string
	var err error

	if !reflect.ValueOf(m.Entity0).IsNil() {
		pk, sk, err = m.Entity0.Keys(0)
		if err != nil {
			return false, err
		}
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
		e0pk, e0sk, err = m.ExtractE0Keys()
		if err != nil {
			return false, err
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
