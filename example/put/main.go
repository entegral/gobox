package main

import (
	"context"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	user, err := exampleLib.PutUser("test@gmail.com", "Test User Name", 30)
	if err != nil {
		panic(err)
	}
	if user != nil {
		logrus.WithField("user", user).Info("User created")
	}
	// put the contact info now:

	contact := &exampleLib.ContactInfo{
		MonoLink: dynamo.NewMonoLink[*exampleLib.User](user),
		Phone:    "555-555-5555",
		Addr:     "123 Main St",
	}
	ctx := context.Background()
	err = contact.Put(ctx, contact)
	if err != nil {
		logrus.Errorln("Error putting contact info")
		panic(err)
	}
	logrus.WithField("contact", contact).Info("Contact info created")
}
