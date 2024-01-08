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
	const year = 2022
	preClearedCar := &exampleLib.Car{
		Make:  carmake,
		Model: model,
		Year:  year,
		Details: &exampleLib.CarDetails{
			"doors": 4,
			"color": "black",
		},
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
		t.Run("LoadEntities", func(t *testing.T) {
			t.Run("LoadEntity0", func(t *testing.T) {
				// check the pink slip only using entities with composite key values
				minimalUser := &exampleLib.User{
					Email: email,
				}
				minimalCar := &exampleLib.Car{
					Make:  carmake,
					Model: model,
				}
				pinkSlip := &exampleLib.PinkSlip{
					DiLink: *dynamo.NewDiLink(minimalUser, minimalCar),
				}
				// now the pink slip has been created locally with only enough info to generate the composite keys
				loaded, err := pinkSlip.LoadEntity0(ctx)
				assert.Equal(t, nil, err)
				assert.Equal(t, true, loaded)
				assert.Equal(t, email, pinkSlip.DiLink.Entity0.Email)
				assert.Equal(t, name, pinkSlip.DiLink.Entity0.Name)
				assert.Equal(t, age, pinkSlip.DiLink.Entity0.Age)
			})
			t.Run("LoadEntity1", func(t *testing.T) {})
			t.Run("LoadEntities", func(t *testing.T) {})
		})
	})
}
