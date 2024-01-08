package main

import (
	"context"

	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()

	user := exampleLib.PutUser(context.Background(), "")
	logrus.WithField("user", user).Info("User created")

	// put less frequently used contact info into a monolink:
	contact := &exampleLib.ContactInfo{
		Phone: "555-555-5555",
		Addr:  "123 Main St",
	}
	isValid, err := contact.CheckLink(ctx, contact, user)
	if !isValid {
		logrus.Errorln("Error checking contact info")
		panic(err)
	}
	err = contact.Put(ctx, contact)
	if err != nil {
		logrus.Errorln("Error putting contact info")
		panic(err)
	}
	logrus.WithField("contact", contact).Info("Contact info created")
}
