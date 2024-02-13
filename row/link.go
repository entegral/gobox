package row

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/clients"
)

// GenerateEntityKeys generates keys for each entity in the link. It iterates over the entities,
// calls the GenerateKeys method for each entity, and stores the generated keys in a slice.
// If an error occurs during key generation for any entity, it returns the keys generated so far and the error.
func (l Link[LinkType, T]) GenerateEntityKeys(ctx context.Context) ([]Key, error) {
	var keys = make([]Key, len(l.Entities))
	for i, row := range l.Entities {
		err := row.GenerateKeys(ctx)
		if err != nil {
			return keys, err
		}
		keys[i] = Key{Pk: row.Pk, Sk: row.Sk, Index: i, IsEntity: true}
	}
	return keys, nil
}

// Link is a struct that represents a link between multiple entities
// It is used to represent a many-to-many relationship between entities
// The key of the link is the combination of the keys of the entities, applied in order
// The link can be queried on either the entityKeys-index or the entityKeys-index
// The partition keys of the primary composite key, and both GSIs, contain type
// information. This is used to make it easier to query for all links of a certain type.
type Link[LinkType, T Rowable] struct {
	Row[LinkType]

	// EntityKeys is a slice of keys of the entities in the link
	EntityKeys []Key

	// Entities is a slice of entities in the link
	Entities []Row[T]
}

// GenerateKeys generates keys for the link and its entities. It first generates keys for each entity
// by calling the GenerateEntityKeys method. Then, it generates the primary key for the link itself,
// which is a combination of the keys of the entities. It also wraps the primary key of each entity
// with the type of the link. If an error occurs at any point, it returns the error.
func (l Link[T0, T1]) GenerateKeys(ctx context.Context) ([]Key, error) {
	keys, err := l.GenerateEntityKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error generating link entity keys for link row of type: %s. Error: %w", l.Row.Type(), err)
	}

	// we need to do two things with the keys:
	// 1. Generate the primary key for this row, which is the combination of the keys of the entities, and
	// 2. Generate the keys for the entities themselves and store them in an entity GSI
	// The primary key for this row is the combination of the keys of the entities, applied in order

	var pk, sk string
	// Generate the primary key for this row
	for _, key := range keys {
		pkSeg, err := addKeySegment(rowPk, key.Pk)
		if err != nil {
			return nil, fmt.Errorf("Error adding key segment to primary key for link row of type: %s. Error: %w", l.Row.Type(), err)
		}
		pk += pkSeg

		// also wrap the entity Pk with the type of the link
		key.Pk, err = prependWithRowType(&l.Row, key.Pk)
		if err != nil {
			return nil, fmt.Errorf("Error prepending row type to entity Pk for link row of type: %s. Error: %w", l.Row.Type(), err)
		}

		skSeg, err := addKeySegment(rowSk, key.Sk)
		if err != nil {
			return nil, fmt.Errorf("Error adding key segment to primary key for link row of type: %s. Error: %w", l.Row.Type(), err)
		}
		sk += skSeg
	}

	// prefix the primary key with the type of the link
	l.Pk, err = prependWithRowType(&l.Row, l.Pk)
	return keys, nil
}

// LoadFromKey loads a Link from DynamoDB using a provided key. It determines which index to use based on the key's Index field,
// wraps the key with the type of the link, and performs a GetItem or Query operation on DynamoDB. If the key's Index is 0,
// it performs a GetItem operation. If the key's Index is not 0, it performs a Query operation on the appropriate GSI.
// After the operation, it unmarshals the returned item(s) into the Link. If multiple items are returned from the Query operation,
// it unmarshals the first item into the Link, and returns the keys of the remaining items in the additionalItems slice.
// If an error occurs at any point, it returns the error.
func (l Link[T0, T1]) LoadFromKey(ctx context.Context, key Key) (additionalItems []Link[T0, T1], nextPage *Key, err error) {
	// wrap the key with the type of the link to make sure only links to this type are returned
	pk, err := prependWithRowType(&l.Row, key.Pk)
	if err != nil {
		return nil, nil, fmt.Errorf("Error prepending row type to entity Pk for link row of type: %s. Error: %w", l.Row.Type(), err)
	}
	client := clients.GetDefaultClient(ctx)
	if key.Index == 0 {
		i := dynamodb.GetItemInput{
			TableName: l.TableName(),
			Key: map[string]awstypes.AttributeValue{
				"pk": &awstypes.AttributeValueMemberS{Value: pk},
				"sk": &awstypes.AttributeValueMemberS{Value: key.Sk},
			},
		}
		out, getErr := client.Dynamo().GetItem(ctx, &i)
		if getErr != nil {
			return nil, nil, fmt.Errorf("Error getting link row of type: %s. Error: %w", l.Row.Type(), getErr)
		}
		unMarshalErr := l.unmarshalMap(out.Item)
		if unMarshalErr != nil {
			return nil, nil, fmt.Errorf("Error unmarshaling link row of type: %s. Error: %w", l.Row.Type(), unMarshalErr)
		}
		return nil, nil, nil
	}
	kce := "pk = :pk AND begins_with(sk, :sk)"
	i := dynamodb.QueryInput{
		TableName:              l.TableName(),
		IndexName:              key.IndexName(),
		KeyConditionExpression: &kce,
		ExpressionAttributeValues: map[string]awstypes.AttributeValue{
			":pk": &awstypes.AttributeValueMemberS{Value: pk},
			":sk": &awstypes.AttributeValueMemberS{Value: key.Sk},
		},
		ExclusiveStartKey: key.LastEvaluatedKey, // Use the LastEvaluatedKey field in case of pagination
	}
	out, err := client.Dynamo().Query(ctx, &i)
	if err != nil {
		return nil, nil, fmt.Errorf("Error querying link row of type: %s. Error: %w", l.Row.Type(), err)
	}
	if len(out.Items) <= 0 {
		return nil, nil, fmt.Errorf("No link row of type: %s found", l.Row.Type())
	}
	if len(out.Items) == 1 {
		err = l.unmarshalMap(out.Items[0])
		if err != nil {
			return nil, nil, fmt.Errorf("Error unmarshaling link row of type: %s. Error: %w", l.Row.Type(), err)
		}
	}
	restOfItems := out.Items[1:]
	err = attributevalue.UnmarshalListOfMaps(restOfItems, &additionalItems)
	if err != nil {
		return nil, nil, fmt.Errorf("Error unmarshaling additional link rows of type: %s. Error: %w", l.Row.Type(), err)
	}
	if out.LastEvaluatedKey != nil {
		key.LastEvaluatedKey = out.LastEvaluatedKey
		return additionalItems, &key, nil
	}
	return additionalItems, nil, nil
}

// NewLink creates a new Link with the provided link type and entities. It creates a new Row with the link type,
// and assigns the entities to the Entities field of the Link.
func NewLink[LinkType, T Rowable](link LinkType, entities []Row[T]) Link[LinkType, T] {
	row := NewRow(link)
	return Link[LinkType, T]{
		Row:      row,
		Entities: entities,
	}
}
