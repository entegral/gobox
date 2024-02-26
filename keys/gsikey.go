package keys

import (
	"fmt"
	"strconv"
)

const (
	IndexFormat = "%s-%s-index"
	PkKey       = "pk"
	SkKey       = "sk"
	Pk1Key      = "pk1"
	Sk1Key      = "sk1"
	Pk2Key      = "pk2"
	Sk2Key      = "sk2"
	Pk3Key      = "pk3"
	Sk3Key      = "sk3"
	Pk4Key      = "pk4"
	Sk4Key      = "sk4"
	Pk5Key      = "pk5"
	Sk5Key      = "sk5"
	Pk6Key      = "pk6"
	Sk6Key      = "sk6"
)

type GSI struct {
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
}

func (g *GSI) GetKeyForGSI(gsi int) (Key, error) {
	pkKey := "pk" + strconv.Itoa(gsi)
	skKey := "sk" + strconv.Itoa(gsi)
	pk := ""
	sk := ""
	switch gsi {
	case 1:
		if g.Pk1 == nil || g.Sk1 == nil || *g.Pk1 == "" || *g.Sk1 == "" {
			return Key{}, fmt.Errorf("GSI1 is not set")
		}
		pk = *g.Pk1
		sk = *g.Sk1
	case 2:
		if g.Pk2 == nil || g.Sk2 == nil || *g.Pk2 == "" || *g.Sk2 == "" {
			return Key{}, fmt.Errorf("GSI2 is not set")
		}
		pk = *g.Pk2
		sk = *g.Sk2
	case 3:
		if g.Pk3 == nil || g.Sk3 == nil || *g.Pk3 == "" || *g.Sk3 == "" {
			return Key{}, fmt.Errorf("GSI3 is not set")
		}
		pk = *g.Pk3
		sk = *g.Sk3
	case 4:
		if g.Pk4 == nil || g.Sk4 == nil || *g.Pk4 == "" || *g.Sk4 == "" {
			return Key{}, fmt.Errorf("GSI4 is not set")
		}
		pk = *g.Pk4
		sk = *g.Sk4
	case 5:
		if g.Pk5 == nil || g.Sk5 == nil || *g.Pk5 == "" || *g.Sk5 == "" {
			return Key{}, fmt.Errorf("GSI5 is not set")
		}
		pk = *g.Pk5
		sk = *g.Sk5
	case 6:
		if g.Pk6 == nil || g.Sk6 == nil || *g.Pk6 == "" || *g.Sk6 == "" {
			return Key{}, fmt.Errorf("GSI6 is not set")
		}
		pk = *g.Pk6
		sk = *g.Sk6
	default:
		return Key{}, fmt.Errorf("GSI %d is not supported", gsi)
	}
	return Key{
		pkKey:        pkKey,
		skKey:        skKey,
		PartitionKey: pk,
		SortKey:      sk,
	}, nil
}
