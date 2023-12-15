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

	// now that we have the user, get can get the contact info in
	// one of two ways:
	contactInfo := &exampleLib.ContactInfo{
		MonoLink: dynamo.NewMonoLink(user),
	}
	loaded, err = contactInfo.CheckLink(ctx, contactInfo, user)
	if err != nil {
		panic(err)
	}
	if loaded {
		logrus.Println("----------------")
		logrus.WithField("contact", *contactInfo).Info("Contact found")
	} else {
		logrus.Info("Contact not found")
	}
}
