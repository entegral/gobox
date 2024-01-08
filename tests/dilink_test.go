package tests

import (
	"context"
	"os"
	"testing"

	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/stretchr/testify/assert"
)

func TestDiLink(t *testing.T) {
	// init
	t.Setenv("TABLENAME", os.Getenv("TABLENAME"))
	t.Setenv("TESTING", os.Getenv("TESTING")) // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	const email = "testDiLinkEmail@gmail.com"
	const name = "TestDiLinkName"
	const age = 33
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
	const carmake = "TestDiLinkMake"
	const model = "TestDiLinkModel"
	preClearCar := &exampleLib.Car{
		Make:  carmake,
		Model: model,
	}
	// ensure we start with an empty slate
	err = preClearCar.Delete(ctx, preClearCar)
	if err != nil {
		t.Error(err)
	}

	t.Run("DiLink", func(t *testing.T) {
		t.Run("CheckLink", func(t *testing.T) {
			missingUser := exampleLib.User{Email: "wontbe@here.com"}
			missingUser.Delete(ctx, &missingUser)
			missingCar := exampleLib.Car{Make: "wontbehere", Model: "wontbehere"}
			missingCar.Delete(ctx, &missingCar)
			t.Run("Should return entity0 not found error", func(t *testing.T) {
				pinkSlip := &exampleLib.PinkSlip{}
				_, err := pinkSlip.CheckLink(ctx, pinkSlip, &missingUser, &missingCar)
				if err == nil {
					t.Error("expected error")
				}
				assert.Equal(t, "user not found", err.Error())
			})
			t.Run("Should return entity1 not found error", func(t *testing.T) {
				err := missingUser.Put(ctx, &missingUser)
				if err != nil {
					t.Error(err)
				}
				pinkSlip := &exampleLib.PinkSlip{}
				_, err = pinkSlip.CheckLink(ctx, pinkSlip, &missingUser, &missingCar)
				if err == nil {
					t.Error("expected error")
				}
				assert.Equal(t, "car not found", err.Error())
			})
			t.Run("Should not return an error", func(t *testing.T) {
				// t.
				err := missingCar.Put(ctx, &missingCar)
				if err != nil {
					t.Error(err)
				}
				// delete to ensure clean state for test
				pinkSlip := &exampleLib.PinkSlip{}
				_, _ = pinkSlip.CheckLink(ctx, pinkSlip, &missingUser, &missingCar)
				err = pinkSlip.Delete(ctx, pinkSlip)
				if err != nil {
					t.Error(err)
				}
				t.Run("Should return false when the record isn't in dynamo", func(t *testing.T) {
					linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, &missingUser, &missingCar)
					assert.Equal(t, nil, err)
					assert.Equal(t, false, linkExists)
				})
				t.Run("should return true when the record is in dynamo", func(t *testing.T) {
					err = pinkSlip.Put(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, &missingUser, &missingCar)
					assert.Equal(t, nil, err)
					assert.Equal(t, true, linkExists)
				})
			})
		})
	})
}
