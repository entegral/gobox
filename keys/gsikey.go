package keys

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
	Pk1 *string `json:"pk1,omitempty"`
	Sk1 *string `json:"sk1,omitempty"`

	Pk2 *string `json:"pk2,omitempty"`
	Sk2 *string `json:"sk2,omitempty"`

	Pk3 *string `json:"pk3,omitempty"`
	Sk3 *string `json:"sk3,omitempty"`

	Pk4 *string `json:"pk4,omitempty"`
	Sk4 *string `json:"sk4,omitempty"`

	Pk5 *string `json:"pk5,omitempty"`
	Sk5 *string `json:"sk5,omitempty"`

	Pk6 *string `json:"pk6,omitempty"`
	Sk6 *string `json:"sk6,omitempty"`
}
