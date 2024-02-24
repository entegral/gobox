package keys

import "fmt"

type Key struct {
	pkKey        string
	skKey        string
	PartitionKey string `dynamodbav:"pk" json:"pk"`
	SortKey      string `dynamodbav:"sk" json:"sk"`
}

func (k *Key) SetIndex(pkKey, skKey string) {
	k.pkKey = pkKey
	k.skKey = skKey
}

func (k *Key) IndexName() *string {
	if k.pkKey == "" && k.skKey == "" {
		return nil
	}
	s := fmt.Sprintf("%s-%s-index", k.pkKey, k.skKey)
	return &s
}

func NewGSIKey(pkKey, skKey, partitionKey, sortKey string) Key {
	return Key{
		pkKey:        pkKey,
		skKey:        skKey,
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}
}

func NewPrimaryKey(partitionKey, sortKey string) Key {
	return Key{
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}
}
