package main

import (
	"context"

	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/entegral/gobox/message"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	user := exampleLib.User{
		Email: "test@gmail.com",
		Name:  "Test User",
		Age:   30,
	}
	ctx := context.Background()
	out, err := message.Send(ctx, "test-queue", user)
	if err != nil {
		logrus.WithError(err).Errorln("error sending message")
		return
	}
	logrus.WithField("MessageId", *out.MessageId).Println("successfully sent message")
}
