package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx := context.Background()

	// put the user and car first:
	users, cars := exampleLib.SeedUsers(ctx, 2), exampleLib.SeedCars(ctx, 1)
	buyer := exampleLib.Buyer{User: users[0]}
	seller := exampleLib.Seller{User: users[1]}
	sale := &exampleLib.Sale{
		TriLink: dynamo.NewTriLink(&buyer, cars[0], &seller),
		Date:    "2018-01-01",
	}
	err := sale.Put(ctx, sale)
	if err != nil {
		logrus.WithError(err).Println("error putting sale")
		panic(err)
	}
	btes, _ := json.MarshalIndent(sale, "", "  ")
	logrus.Info("sale created")
	fmt.Printf("%s", string(btes))
}
