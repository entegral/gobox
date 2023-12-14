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
	out, err := message.BroadcastMessage(ctx, "newUser", user)
	if err != nil {
		panic(err)
	}
	if len(out.Entries) == 0 {
		panic("no entries returned")
	}
	for _, entry := range out.Entries {
		if entry.ErrorCode != nil {
			logrus.WithField("Code", *entry.ErrorCode).Printf("error: %s", *entry.ErrorMessage)
			continue
		}
		logrus.WithField("Entry", entry).Println("successfully broadcasted message")
	}
}
