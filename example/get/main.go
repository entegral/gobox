package main

import (
	"context"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()
	user := &exampleLib.User{
		Email: "test@gmail.com",
	}
	loaded, err := user.Get(ctx, user)
	if err != nil {
		panic(err)
	}
	if loaded {
		logrus.WithField("user", *user).Info("User found")
	} else {
		logrus.Info("User not found")
	}

	// now get the contact info:
	contact := &exampleLib.ContactInfo{
		MonoLink: *dynamo.NewMonoLink[*exampleLib.User](user),
	}
	loaded, err = contact.Get(ctx, contact)
	if err != nil {
		panic(err)
	}
	if loaded {
		logrus.Println("----------------")
		logrus.WithField("contact", *contact).Info("Contact found")
	} else {
		logrus.Info("Contact not found")
	}
}
