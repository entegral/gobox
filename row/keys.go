package row

type Keys struct {
	Pk  string `dynamodbav:"pk,omitempty" json:"pk,omitempty"`
	Sk  string `dynamodbav:"sk,omitempty" json:"sk,omitempty"`
	Pk1 string `dynamodbav:"pk1,omitempty" json:"pk1,omitempty"`
	Sk1 string `dynamodbav:"sk1,omitempty" json:"sk1,omitempty"`
	Pk2 string `dynamodbav:"pk2,omitempty" json:"pk2,omitempty"`
	Sk2 string `dynamodbav:"sk2,omitempty" json:"sk2,omitempty"`
	Pk3 string `dynamodbav:"pk3,omitempty" json:"pk3,omitempty"`
	Sk3 string `dynamodbav:"sk3,omitempty" json:"sk3,omitempty"`
	Pk4 string `dynamodbav:"pk4,omitempty" json:"pk4,omitempty"`
	Sk4 string `dynamodbav:"sk4,omitempty" json:"sk4,omitempty"`
	Pk5 string `dynamodbav:"pk5,omitempty" json:"pk5,omitempty"`
	Sk5 string `dynamodbav:"sk5,omitempty" json:"sk5,omitempty"`
	Pk6 string `dynamodbav:"pk6,omitempty" json:"pk6,omitempty"`
	Sk6 string `dynamodbav:"sk6,omitempty" json:"sk6,omitempty"`
}

func (k *Keys) MaxGSIs() int {
	return 6
}

// SetKey
func (k *Keys) SetKey(key Key) {
	switch key.Index {
	case 0:
		k.Pk = key.PK
		k.Sk = key.SK
	case 1:
		k.Pk1 = key.PK
		k.Sk1 = key.SK
	case 2:
		k.Pk2 = key.PK
		k.Sk2 = key.SK
	case 3:
		k.Pk3 = key.PK
		k.Sk3 = key.SK
	case 4:
		k.Pk4 = key.PK
		k.Sk4 = key.SK
	case 5:
		k.Pk5 = key.PK
		k.Sk5 = key.SK
	case 6:
		k.Pk6 = key.PK
		k.Sk6 = key.SK
	}
}
