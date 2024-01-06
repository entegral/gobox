package exampleLib

import "context"

func PutCar(ctx context.Context) *Car {
	car := &Car{
		Make:  "Honda",
		Model: "Civic",
		Year:  2018,
	}
	err := car.Put(ctx, car)
	if err != nil {
		panic(err)
	}
	return car
}
