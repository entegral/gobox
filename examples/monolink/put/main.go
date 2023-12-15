package main

import (
	"context"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	email, name, age := "test@gmail.com", "Test User Name", 30
	ctx := context.Background()
	// Create a new user
	user := &exampleLib.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	err := user.Put(ctx, user)
	if err != nil {
		panic(err)
	}
	userWithOtherTable := &exampleLib.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	userWithOtherTable.Tablename = "otherTableName"
	err = userWithOtherTable.Put(ctx, userWithOtherTable)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	if user != nil {
		logrus.WithField("user", user).Info("User created")
	}

	// put the contact info link now:
	contact := &exampleLib.ContactInfo{
		MonoLink: dynamo.NewMonoLink[*exampleLib.User](user),
		Phone:    "555-555-5555",
		Addr:     "123 Main St",
	}

	err = contact.Put(ctx, contact)
	if err != nil {
		logrus.Errorln("Error putting contact info")
		panic(err)
	}
	logrus.WithField("contact", contact).Info("Contact info created")
}
