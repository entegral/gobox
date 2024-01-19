package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/clients"
	ttypes "github.com/entegral/gobox/types"
)

// findLinkRowsByEntityGSI is a generic method to query for a list of rows based on the Entity1.
func findLinkRowsByEntityGSI[T ttypes.Linkable](ctx context.Context, clients *clients.Client, entity T, entityGSI EntityGSI, linkType string) ([]map[string]types.AttributeValue, error) {
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
	tn := entity.TableName(ctx)
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
			":type": &types.AttributeValueMemberS{Value: linkType},
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

// FindLinksByEntity0 is a generic method to query for a list of links based on the Entity0.
func FindLinksByEntity0[T0, CustomLinkType ttypes.Linkable](ctx context.Context, e0 T0, linkWrapper ttypes.Typeable) ([]CustomLinkType, error) {
	client := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntityGSI[T0](ctx, client, e0, Entity0GSI, linkWrapper.Type())
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

// FindLinksByEntity1 is a generic method to query for a list of links based on the Entity1.
func FindLinksByEntity1[T1, CustomLinkType ttypes.Linkable](ctx context.Context, e1 T1, linkType string) ([]CustomLinkType, error) {
	clients := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntityGSI[T1](ctx, clients, e1, Entity1GSI, linkType)
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
	return links, err
}

// FindLinksByEntity2 is a generic method to query for a list of links based on the Entity2.
func FindLinksByEntity2[T2, CustomLinkType ttypes.Linkable](ctx context.Context, e2 T2, linkWrapper ttypes.Typeable) ([]CustomLinkType, error) {
	clients := clients.GetDefaultClient(ctx)
	rows, err := findLinkRowsByEntityGSI[T2](ctx, clients, e2, Entity2GSI, linkWrapper.Type())
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
