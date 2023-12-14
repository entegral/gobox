package main

import (
	"context"

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
	contactInfo := &exampleLib.ContactInfo{}
	loaded, err = contactInfo.GetFrom(contactInfo, user)
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
