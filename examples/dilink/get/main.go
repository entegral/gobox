package main

import (
	"context"

	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()
	user := exampleLib.PutUser(ctx, "")
	car := exampleLib.PutCar(ctx, "", "")
	logrus.WithFields(logrus.Fields{
		"user": user,
		"car":  car,
	}).Info("data seeded")

	// get the pinkSlip link now:
	pinkSlip := &exampleLib.PinkSlip{}
	loaded, err := pinkSlip.CheckLink(ctx, pinkSlip, user, car)
	if err != nil {
		panic(err)
	}
	if !loaded {
		logrus.Errorln("Pink slip not found")
		return
	}
	logrus.WithField("pink slip", pinkSlip).Info("Pink slip loaded")
}
