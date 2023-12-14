package exampleLib

import (
	"context"
	"os"
)

func PutUser(email, name string, age int) (*User, error) {
	// ensure env var TABLENAME is set
	if tn := os.Getenv("TABLENAME"); tn == "" {
		panic("TABLENAME env var not set")
	}
	ctx := context.Background()
	// Create a new user
	user := &User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	err := user.Put(ctx, user)
	if err != nil {
		panic(err)
	}
	return user, nil
}
