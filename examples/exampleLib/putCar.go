package exampleLib

import "context"

func PutCar(ctx context.Context, make, model string) *Car {
	if make == "" {
		make = "Honda"
	}
	if model == "" {
		model = "Civic"
	}
	car := &Car{
		Make:  make,
		Model: model,
		Year:  2018,
	}
	err := car.Put(ctx, car)
	if err != nil {
		panic(err)
	}
	return car
}
