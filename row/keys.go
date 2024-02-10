package row

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dgryski/trifles/uuid"
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
	Keys(gsi int) (pk, sk string, err error)
}

type Keys struct {
	Pk string `dynamodbav:"pk" json:"pk"`
	Sk string `dynamodbav:"sk" json:"sk"`

	Pk1 *string `dynamodbav:"pk1,omitempty" json:"pk1,omitempty"`
	Sk1 *string `dynamodbav:"sk1,omitempty" json:"sk1,omitempty"`

	Pk2 *string `dynamodbav:"pk2,omitempty" json:"pk2,omitempty"`
	Sk2 *string `dynamodbav:"sk2,omitempty" json:"sk2,omitempty"`

	Pk3 *string `dynamodbav:"pk3,omitempty" json:"pk3,omitempty"`
	Sk3 *string `dynamodbav:"sk3,omitempty" json:"sk3,omitempty"`

	Pk4 *string `dynamodbav:"pk4,omitempty" json:"pk4,omitempty"`
	Sk4 *string `dynamodbav:"sk4,omitempty" json:"sk4,omitempty"`

	Pk5 *string `dynamodbav:"pk5,omitempty" json:"pk5,omitempty"`
	Sk5 *string `dynamodbav:"sk5,omitempty" json:"sk5,omitempty"`

	Pk6 *string `dynamodbav:"pk6,omitempty" json:"pk6,omitempty"`
	Sk6 *string `dynamodbav:"sk6,omitempty" json:"sk6,omitempty"`

	// KeyProcessor is a function that processes the keys after they are generated.
	// This is useful for adding prefixes or suffixes to the keys.
	KeyProcessor func(item Rowable, pk, sk string) (string, string, error)
}

// AddKeysToMap adds the keys to the dynamodb attribute map.
func (k *Keys) AddKeysToMap(m map[string]types.AttributeValue) {
	m["pk"] = &types.AttributeValueMemberS{Value: k.Pk}
	m["sk"] = &types.AttributeValueMemberS{Value: k.Sk}
	if k.Pk1 != nil && k.Sk1 != nil {
		m["pk1"] = &types.AttributeValueMemberS{Value: *k.Pk1}
		m["sk1"] = &types.AttributeValueMemberS{Value: *k.Sk1}
	}
	if k.Pk2 != nil && k.Sk2 != nil {
		m["pk2"] = &types.AttributeValueMemberS{Value: *k.Pk2}
		m["sk2"] = &types.AttributeValueMemberS{Value: *k.Sk2}
	}
	if k.Pk3 != nil && k.Sk3 != nil {
		m["pk3"] = &types.AttributeValueMemberS{Value: *k.Pk3}
		m["sk3"] = &types.AttributeValueMemberS{Value: *k.Sk3}
	}
	if k.Pk4 != nil && k.Sk4 != nil {
		m["pk4"] = &types.AttributeValueMemberS{Value: *k.Pk4}
		m["sk4"] = &types.AttributeValueMemberS{Value: *k.Sk4}
	}
	if k.Pk5 != nil && k.Sk5 != nil {
		m["pk5"] = &types.AttributeValueMemberS{Value: *k.Pk5}
		m["sk5"] = &types.AttributeValueMemberS{Value: *k.Sk5}
	}
	if k.Pk6 != nil && k.Sk6 != nil {
		m["pk6"] = &types.AttributeValueMemberS{Value: *k.Pk6}
		m["sk6"] = &types.AttributeValueMemberS{Value: *k.Sk6}
	}
}

func (k *Keys) ProcessKeys(item Rowable, pk, sk string) (string, string, error) {
	if k.KeyProcessor != nil {
		return k.KeyProcessor(item, pk, sk)
	}
	pk, err := prependWithRowType(item, pk)
	if err != nil {
		return "", "", err
	}

	sk, err = addKeySegment(rowSk, sk)
	if err != nil {
		return "", "", err
	}
	return pk, sk, nil
}

// GenerateKeys accepts any Keyable item and generates the keys for the row.
// If the keys are already set, then it overwrites them.
func (k *Keys) GenerateKeys(item Rowable) error {
	for i := 0; i < k.MaxGSIs(); i++ {
		pk, sk, err := item.Keys(i)
		if err != nil {
			return err
		}
		if pk != "" && sk != "" {

			pk, sk, err := k.ProcessKeys(item, pk, sk)
			if err != nil {
				return err
			}

			switch i {
			case 0:
				k.Pk = pk
				k.Sk = sk
			case 1:
				k.Pk1 = &pk
				k.Sk1 = &sk
			case 2:
				k.Pk2 = &pk
				k.Sk2 = &sk
			case 3:
				k.Pk3 = &pk
				k.Sk3 = &sk
			case 4:
				k.Pk4 = &pk
				k.Sk4 = &sk
			case 5:
				k.Pk5 = &pk
				k.Sk5 = &sk
			case 6:
				k.Pk6 = &pk
				k.Sk6 = &sk
			}
		} else if pk == "" && sk == "" && i == 0 {
			k.Pk = uuid.UUIDv4()
			k.Sk = "default"
			continue
		} else if pk == "" && sk == "" {
			continue
		} else if pk == "" {
			return ErrPartitionKeyRequired{GSI: i, Type: item.Type()}
		} else if sk == "" {
			return ErrSortKeyRequired{GSI: i, Type: item.Type()}
		}
	}
	return nil
}

type ErrPartitionKeyRequired struct {
	GSI  int
	Type string
}

func (e ErrPartitionKeyRequired) Error() string {
	return fmt.Sprintf("partition key is required for gsi %d of type %s", e.GSI, e.Type)
}

type ErrSortKeyRequired struct {
	GSI  int
	Type string
}

func (e ErrSortKeyRequired) Error() string {
	return fmt.Sprintf("sort key is required for gsi %d of type %s", e.GSI, e.Type)
}

func (k *Keys) MaxGSIs() int {
	return 6
}
