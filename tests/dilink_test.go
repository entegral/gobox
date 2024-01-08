package tests

import (
	"context"
	"os"
	"testing"

	"github.com/entegral/gobox/examples/exampleLib"
)

func TestDiLink(t *testing.T) {
	// init
	os.Setenv("TABLENAME", "arctica")
	os.Setenv("TESTING", "true") // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	email := "testDiLinkEmail@gmail.com"
	name := "TestDiLinkName"
	age := 33
	preClearUser := &exampleLib.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	// ensure we start with an empty slate
	err := preClearUser.Delete(ctx, preClearUser)
	if err != nil {
		t.Error(err)
	}

	preClearCar := &exampleLib.Car{
		Make:  "TestDiLinkMake",
		Model: "TestDiLinkModel",
	}
	// ensure we start with an empty slate
	err = preClearCar.Delete(ctx, preClearCar)
	if err != nil {
		t.Error(err)
	}

	t.Run("DiLink", func(t *testing.T) {

	})
}
