package tests

import (
	"context"
	"os"
	"testing"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/stretchr/testify/assert"
)

func TestDiLink(t *testing.T) {
	// init
	t.Setenv("TABLENAME", os.Getenv("TABLENAME"))
	t.Setenv("TESTING", os.Getenv("TESTING")) // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	// ensure we start with an empty slate
	const email = "testDiLinkEmail@gmail.com"
	const name = "TestDiLinkName"
	const age = 33
	preClearedUser := &exampleLib.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	err := preClearedUser.Delete(ctx, preClearedUser)
	if err != nil {
		t.Error(err)
	}
	// ensure we start with an empty slate
	const carmake = "TestDiLinkMake"
	const model = "TestDiLinkModel"
	preClearedCar := &exampleLib.Car{
		Make:  carmake,
		Model: model,
	}
	err = preClearedCar.Delete(ctx, preClearedCar)
	if err != nil {
		t.Error(err)
	}
	// ensure we start with an empty slate
	pinkSlip := &exampleLib.PinkSlip{
		DiLink: *dynamo.NewDiLink(preClearedUser, preClearedCar),
	}
	err = pinkSlip.Delete(ctx, pinkSlip)
	if err != nil {
		t.Error(err)
	}
	t.Run("DiLink", func(t *testing.T) {
		t.Run("CheckLink", func(t *testing.T) {
			t.Run("Should return an error", func(t *testing.T) {
				t.Run("entity0 not found error", func(t *testing.T) {
					_, err := pinkSlip.CheckLink(ctx, pinkSlip, preClearedUser, preClearedCar)
					if err == nil {
						t.Error("expected error")
					}
					assert.Equal(t, "user not found", err.Error())
				})
				t.Run("entity1 not found error", func(t *testing.T) {
					// create the user now
					err := preClearedUser.Put(ctx, preClearedUser)
					if err != nil {
						t.Error(err)
					}
					_, err = pinkSlip.CheckLink(ctx, pinkSlip, preClearedUser, preClearedCar)
					if err == nil {
						t.Error("expected error")
					}
					assert.Equal(t, "car not found", err.Error())
				})
			})
			t.Run("Should not return an error", func(t *testing.T) {
				// create the car now
				err := preClearedCar.Put(ctx, preClearedCar)
				if err != nil {
					t.Error(err)
				}
				t.Run("Should return false when the record isn't in dynamo", func(t *testing.T) {
					linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, preClearedUser, preClearedCar)
					assert.Equal(t, nil, err)
					assert.Equal(t, false, linkExists)
				})
				t.Run("should return true when the record is in dynamo", func(t *testing.T) {
					// finally create the pink slip
					err = pinkSlip.Put(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, preClearedUser, preClearedCar)
					assert.Equal(t, nil, err)
					assert.Equal(t, true, linkExists)
				})
			})
		})
	})
}
