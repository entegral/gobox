package keys

type PrimaryKey struct {
	PartitionKey string `json:"pk"`
	SortKey      string `json:"sk"`
}

type GSI1 struct {
	Pk1 string `json:"pk1"`
	Sk1 string `json:"sk1"`
}

type GSI2 struct {
	Pk2 string `json:"pk2"`
	Sk2 string `json:"sk2"`
}

type GSI3 struct {
	Pk3 string `json:"pk3"`
	Sk3 string `json:"sk3"`
}

type GSI4 struct {
	Pk4 string `json:"pk4"`
	Sk4 string `json:"sk4"`
}

type GSI5 struct {
	Pk5 string `json:"pk5"`
	Sk5 string `json:"sk5"`
}

type GSI6 struct {
	Pk6 string `json:"pk6"`
	Sk6 string `json:"sk6"`
}
