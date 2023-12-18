package exampleLib

import "github.com/entegral/gobox/dynamo"

type PinkSlip struct {
	*dynamo.DiLink[*User, *Car]
	DateOfPurchase string
}
