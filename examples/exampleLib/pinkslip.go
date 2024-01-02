package exampleLib

import "github.com/entegral/gobox/dynamo"

// PinkSlip is a link between a user and a car, along with representing the date of purchase
// and the vehicles specific VIN.
// It is a DiLink, meaning it is a link between two entities, and it is bidirectional.
//
// Since it is a DiLink, it uses the Linkable implementations of User and Car to
// deterministically generate the primary key and sort key for the link. This maintains
// the entropy of the link, and allows for easy querying of the link.
type PinkSlip struct {
	*dynamo.DiLink[*User, *Car]
	DateOfPurchase string
	VIN            string
}

// but what if you need to add another relation to the pinkslip? for example by
// linking to a dealership? you can do that by
