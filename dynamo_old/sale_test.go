package dynamo

type Buyer struct {
	*User
}

func (b *Buyer) Type() string {
	return "buyer"
}

type Seller struct {
	*User
}

func (s *Seller) Type() string {
	return "seller"
}

type Sale struct {
	*TriLink[*Buyer, *Car, *Seller]
	// the old owner of the car
	Date string
}
