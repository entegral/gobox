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
	user, car := seedRows(ctx)

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

func seedRows(ctx context.Context) (user *exampleLib.User, car *exampleLib.Car) {
	email, name, age := "test@gmail.com", "Test User Name", 30
	user = &exampleLib.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	err := user.Put(ctx, user)
	if err != nil {
		panic(err)
	}
	car = &exampleLib.Car{
		Make:  "Honda",
		Model: "Civic",
		Year:  2018,
	}
	err = car.Put(ctx, car)
	if err != nil {
		panic(err)
	}
	return user, car
}
