package exampleLib

import (
	"context"
	"fmt"
)

func SeedUsers(ctx context.Context, count int) []*User {
	var users []*User
	for i := 0; i < count; i++ {
		users = append(users, seedUser(ctx, i))
	}
	return users
}

func seedUser(ctx context.Context, index int) *User {
	email, name, age := "test@gmail.com", "Test User Name", 30
	user := &User{
		Email: fmt.Sprintf("%s:%d", email, index),
		Name:  name,
		Age:   age,
	}
	err := user.Put(ctx, user)
	if err != nil {
		panic(err)
	}
	return user
}

func SeedCars(ctx context.Context, count int) []*Car {
	var cars []*Car
	for i := 0; i < count; i++ {
		cars = append(cars, seedCar(ctx, i))
	}
	return cars
}

func seedCar(ctx context.Context, index int) *Car {
	car := &Car{
		Make:  fmt.Sprintf("Honda:%d", index),
		Model: "Civic",
		Year:  2018,
	}
	err := car.Put(ctx, car)
	if err != nil {
		panic(err)
	}
	return car
}
