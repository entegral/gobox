package row

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Keyable is an interface for types so they can be used with the
// GetItem, PutItem, and DeleteItem functions.
type Keyable interface {

	// Keys returns the partition key and sort key for the given GSI.
	// If the GSI is 0, then the primary composite key is returned/assumed.
	// If the GSI is 1, then the composite key for the pk1-sk1-index is assumed.
	// When implementing this method, you should return the appropriate
	// partition key and sort key for the given GSI, however, you should
	// also ensure any other GSI fields that rely on struct fields are
	// populated as well.
	Keys(gsi int) (key Key, err error)
}

type Keys struct {
	PrimaryKey Key
	GSI1       Key
	GSI2       Key
	GSI3       Key
	GSI4       Key
	GSI5       Key
	GSI6       Key

	// Linked entities
	Entity0 Key
	Entity1 Key
	Entity2 Key
}

func (k *Keys) unmarshalKeysFromMap(m map[string]awstypes.AttributeValue) error {
	for key := range m {
		err := k.unmarshalKeyFromMap(m, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keys) MarshalMap(m map[string]awstypes.AttributeValue) error {
	if k.PrimaryKey.PK != "" && k.PrimaryKey.SK != "" {
		m["pk"] = &types.AttributeValueMemberS{Value: k.PrimaryKey.PK}
		m["sk"] = &types.AttributeValueMemberS{Value: k.PrimaryKey.SK}
	}
	if k.GSI1.PK != "" && k.GSI1.SK != "" {
		m["pk1"] = &types.AttributeValueMemberS{Value: k.GSI1.PK}
		m["sk1"] = &types.AttributeValueMemberS{Value: k.GSI1.SK}
	}
	if k.GSI2.PK != "" && k.GSI2.SK != "" {
		m["pk2"] = &types.AttributeValueMemberS{Value: k.GSI2.PK}
		m["sk2"] = &types.AttributeValueMemberS{Value: k.GSI2.SK}
	}
	if k.GSI3.PK != "" && k.GSI3.SK != "" {
		m["pk3"] = &types.AttributeValueMemberS{Value: k.GSI3.PK}
		m["sk3"] = &types.AttributeValueMemberS{Value: k.GSI3.SK}
	}
	if k.GSI4.PK != "" && k.GSI4.SK != "" {
		m["pk4"] = &types.AttributeValueMemberS{Value: k.GSI4.PK}
		m["sk4"] = &types.AttributeValueMemberS{Value: k.GSI4.SK}
	}
	if k.GSI5.PK != "" && k.GSI5.SK != "" {
		m["pk5"] = &types.AttributeValueMemberS{Value: k.GSI5.PK}
		m["sk5"] = &types.AttributeValueMemberS{Value: k.GSI5.SK}
	}
	if k.GSI6.PK != "" && k.GSI6.SK != "" {
		m["pk6"] = &types.AttributeValueMemberS{Value: k.GSI6.PK}
		m["sk6"] = &types.AttributeValueMemberS{Value: k.GSI6.SK}
	}
	if k.Entity0.PK != "" && k.Entity0.SK != "" {
		m["e0pk"] = &types.AttributeValueMemberS{Value: k.Entity0.PK}
		m["e0sk"] = &types.AttributeValueMemberS{Value: k.Entity0.SK}
	}
	if k.Entity1.PK != "" && k.Entity1.SK != "" {
		m["e1pk"] = &types.AttributeValueMemberS{Value: k.Entity1.PK}
		m["e1sk"] = &types.AttributeValueMemberS{Value: k.Entity1.SK}
	}
	if k.Entity2.PK != "" && k.Entity2.SK != "" {
		m["e2pk"] = &types.AttributeValueMemberS{Value: k.Entity2.PK}
		m["e2sk"] = &types.AttributeValueMemberS{Value: k.Entity2.SK}
	}
	return nil
}

func (k *Keys) unmarshalKeyFromMap(m map[string]awstypes.AttributeValue, key string) error {

	if val, ok := m[key]; ok {
		switch key {
		case "pk":
			err := attributevalue.Unmarshal(val, &k.PrimaryKey.PK)
			if err != nil {
				return err
			}
			k.PrimaryKey.Index = 0
		case "sk":
			err := attributevalue.Unmarshal(val, &k.PrimaryKey.SK)
			if err != nil {
				return err
			}
			k.PrimaryKey.Index = 0
		case "pk1":
			err := attributevalue.Unmarshal(val, &k.GSI1.PK)
			if err != nil {
				return err
			}
			k.GSI1.Index = 1
		case "sk1":
			err := attributevalue.Unmarshal(val, &k.GSI1.SK)
			if err != nil {
				return err
			}
			k.GSI1.Index = 1
		case "pk2":
			err := attributevalue.Unmarshal(val, &k.GSI2.PK)
			if err != nil {
				return err
			}
			k.GSI2.Index = 2
		case "sk2":
			err := attributevalue.Unmarshal(val, &k.GSI2.SK)
			if err != nil {
				return err
			}
			k.GSI2.Index = 2
		case "pk3":
			err := attributevalue.Unmarshal(val, &k.GSI3.PK)
			if err != nil {
				return err
			}
			k.GSI3.Index = 3
		case "sk3":
			err := attributevalue.Unmarshal(val, &k.GSI3.SK)
			if err != nil {
				return err
			}
			k.GSI3.Index = 3
		case "pk4":
			err := attributevalue.Unmarshal(val, &k.GSI4.PK)
			if err != nil {
				return err
			}
			k.GSI4.Index = 4
		case "sk4":
			err := attributevalue.Unmarshal(val, &k.GSI4.SK)
			if err != nil {
				return err
			}
			k.GSI4.Index = 4
		case "pk5":
			err := attributevalue.Unmarshal(val, &k.GSI5.PK)
			if err != nil {
				return err
			}
			k.GSI5.Index = 5
		case "sk5":
			err := attributevalue.Unmarshal(val, &k.GSI5.SK)
			if err != nil {
				return err
			}
			k.GSI5.Index = 5
		case "pk6":
			err := attributevalue.Unmarshal(val, &k.GSI6.PK)
			if err != nil {
				return err
			}
			k.GSI6.Index = 6
		case "sk6":
			err := attributevalue.Unmarshal(val, &k.GSI6.SK)
			if err != nil {
				return err
			}
			k.GSI6.Index = 6
		case "e0pk":
			err := attributevalue.Unmarshal(val, &k.Entity0.PK)
			if err != nil {
				return err
			}
			k.Entity0.Index = 0
			k.Entity0.IsEntity = true
		case "e0sk":
			err := attributevalue.Unmarshal(val, &k.Entity0.SK)
			if err != nil {
				return err
			}
			k.Entity0.Index = 0
			k.Entity0.IsEntity = true
		case "e1pk":
			err := attributevalue.Unmarshal(val, &k.Entity1.PK)
			if err != nil {
				return err
			}
			k.Entity1.Index = 1
			k.Entity1.IsEntity = true
		case "e1sk":
			err := attributevalue.Unmarshal(val, &k.Entity1.SK)
			if err != nil {
				return err
			}
			k.Entity1.Index = 1
			k.Entity1.IsEntity = true
		case "e2pk":
			err := attributevalue.Unmarshal(val, &k.Entity2.PK)
			if err != nil {
				return err
			}
			k.Entity2.Index = 2
			k.Entity2.IsEntity = true
		case "e2sk":
			err := attributevalue.Unmarshal(val, &k.Entity2.SK)
			if err != nil {
				return err
			}
			k.Entity2.Index = 2
			k.Entity2.IsEntity = true
		default:
			return nil
		}
	}
	return nil
}

func (k *Keys) MaxGSIs() int {
	return 6
}

// SetKey
func (k *Keys) SetKey(key Key) {
	switch key.Index {
	case 0:
		k.PrimaryKey = key
	case 1:
		k.GSI1 = key
	case 2:
		k.GSI2 = key
	case 3:
		k.GSI3 = key
	case 4:
		k.GSI4 = key
	case 5:
		k.GSI5 = key
	case 6:
		k.GSI6 = key
	}
}
