package types

type Keyable interface {
	Keys(gsi int) (pk string, sk string, err error)
}
