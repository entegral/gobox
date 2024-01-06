package main

import (
	"context"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()

	// put the user and car first:
	user, car := exampleLib.PutUser(ctx), exampleLib.PutCar(ctx)

	// put the pinkSlip link now:
	pinkSlip := &exampleLib.PinkSlip{
		DiLink:         dynamo.NewDiLink(user, car),
		DateOfPurchase: "2018-01-01",
	}

	err := pinkSlip.Put(ctx, pinkSlip)
	if err != nil {
		panic(err)
	}
	logrus.WithField("pink slip", pinkSlip).Info("Pink slip created")
}
