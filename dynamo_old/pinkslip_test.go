package dynamo

import (
	"context"
)

// PinkSlip is a link between a user and a car, along with representing the date of purchase
// and the vehicles specific VIN.
// It is a DiLink, meaning it is a link between two entities, and it is bidirectional.
//
// Since it is a DiLink, it uses the Linkable implementations of User and Car to
// deterministically generate the primary key and sort key for the link. This maintains
// the entropy of the link, and allows for easy querying of the link.
type PinkSlip struct {
	DiLink[*User, *Car]
	DateOfPurchase string `dynamodbav:"dateOfPurchase,omitempty" json:"dateOfPurchase,omitempty"`
	VIN            string `dynamodbav:"vin,omitempty" json:"vin,omitempty"`
}

func (p *PinkSlip) Type() string {
	return "PinkSlip"
}

// Since the DiLink is a tool to help us build up our directory, lets add some helper methods
// to make it easier to use.

func (p *PinkSlip) User(ctx context.Context) (*User, error) {
	_, err := p.LoadEntity0(ctx)
	if err != nil {
		return nil, err
	}
	return p.Entity0, nil
}

func (p *PinkSlip) Users(ctx context.Context) ([]*User, error) {
	users, err := p.LoadEntity0s(ctx, p)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *PinkSlip) Car(ctx context.Context) (*Car, error) {
	_, err := p.LoadEntity1(ctx)
	if err != nil {
		return nil, err
	}
	return p.Entity1, nil
}

func (p *PinkSlip) Cars(ctx context.Context) ([]*Car, error) {
	cars, err := p.LoadEntity1s(ctx, p)
	if err != nil {
		return nil, err
	}
	return cars, nil
}
