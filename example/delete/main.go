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
	err := user.Delete(ctx, user)
	if err != nil {
		panic(err)
	}
	if user.DeleteItemOutput != nil && user.DeleteItemOutput.Attributes != nil {
		logrus.WithField("delete attributes", user.DeleteItemOutput.Attributes).Info("User deleted")
	} else {
		logrus.Info("User not found")
	}

	// now get the contact info:
	contact := &exampleLib.ContactInfo{
		MonoLink: *dynamo.NewMonoLink[*exampleLib.User](user),
	}
	err = contact.Delete(ctx, contact)
	if err != nil {
		panic(err)
	}
	if contact.DeleteItemOutput != nil && contact.DeleteItemOutput.Attributes != nil {
		logrus.WithField("delete attributes", contact.DeleteItemOutput.Attributes).Info("Contact deleted")
	} else {
		logrus.Info("Contact not found")
	}
}
