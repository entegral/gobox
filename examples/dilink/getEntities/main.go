package main

import (
	"context"
	"log/slog"

	"github.com/entegral/gobox/examples/exampleLib"
)

func main() {
	ctx := context.Background()
	// user := exampleLib.PutUser(ctx, "")
	user := exampleLib.PutUser(ctx, "test2@gmail.com")
	car := exampleLib.PutCar(ctx, "", "")
	car2 := exampleLib.PutCar(ctx, "Polestar", "2")
	// car := &exampleLib.Car{}

	pinkSlip := &exampleLib.PinkSlip{}
	linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, user, car)
	if err != nil {
		slog.Error("error checking link:", err)
	}
	if !linkExists {
		slog.Error("PinkSlip not found")
		err = pinkSlip.Put(ctx, pinkSlip)
		if err != nil {
			slog.Error("error putting pinkSlip:", err)
		}
		// now the pink slip should exist for at least this car and user
		// so we should see one result when we query for it below
	}

	// get the cars:
	cars, err := pinkSlip.Cars(ctx)
	if err != nil {
		slog.Error("error getting cars:", err)
	}
	for _, car := range cars {
		slog.Info("\n")
		slog.Info("---")
		slog.Info("car:", "make:", car.Make)  // Add an empty string as the final argument
		slog.Info("car:", "model", car.Model) // Add an empty string as the final argument
		slog.Info("---")
		slog.Info("\n")
	}

	// second user with same car
	pinkSlip2 := &exampleLib.PinkSlip{}
	linkExists, err = pinkSlip2.CheckLink(ctx, pinkSlip2, user, car2)
	if err != nil {
		slog.Error("error checking link:", err)
	}
	if !linkExists {
		slog.Error("PinkSlip2 not found")
		err = pinkSlip2.Put(ctx, pinkSlip2)
		if err != nil {
			slog.Error("error putting pinkSlip:", err)
		}
		// now the pink slip should exist for this car and both users
		// so we should see two results when we query for it below
	}
	cars, err = pinkSlip2.Cars(ctx)
	if err != nil {
		slog.Error("error getting cars:", err)
	}
	for _, car := range cars {
		slog.Info("\n")
		slog.Info("---")
		slog.Info("car:", "make:", car.Make)  // Add an empty string as the final argument
		slog.Info("car:", "model", car.Model) // Add an empty string as the final argument
		slog.Info("---")
		slog.Info("\n")
	}
}
