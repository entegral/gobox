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
	var carDetails = exampleLib.CarDetails{
		"doors": 4,
		"color": "black",
	}
	preClearedCar := &exampleLib.Car{
		Make:    carmake,
		Model:   model,
		Year:    year,
		Details: &carDetails,
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
			// check the pink slip only using entities with composite key values
			minimalUser := &exampleLib.User{
				Email: email,
			}
			minimalCar := &exampleLib.Car{
				Make:  carmake,
				Model: model,
				Year:  year,
			}
			pinkSlip := &exampleLib.PinkSlip{
				DiLink: *dynamo.NewDiLink(minimalUser, minimalCar),
			}
			t.Run("LoadEntity0", func(t *testing.T) {
				// now the pink slip has been created locally with only enough info to generate the composite keys
				loaded, err := pinkSlip.LoadEntity0(ctx)
				assert.Equal(t, nil, err)
				assert.Equal(t, true, loaded)
				assert.Equal(t, email, pinkSlip.Entity0.Email)
				assert.Equal(t, name, pinkSlip.Entity0.Name)
				assert.Equal(t, age, pinkSlip.Entity0.Age)
			})
			t.Run("LoadEntity1", func(t *testing.T) {
				// ensure car is in dynamo
				err := preClearedCar.Put(ctx, preClearedCar)
				if err != nil {
					t.Error(err)
				}
				// now the pink slip has been created locally with only enough info to generate the composite keys
				loaded, err := pinkSlip.LoadEntity1(ctx)
				assert.Equal(t, nil, err)
				assert.Equal(t, true, loaded)
				assert.Equal(t, carmake, pinkSlip.Entity1.Make)
				assert.Equal(t, model, pinkSlip.Entity1.Model)
				assert.Equal(t, year, pinkSlip.Entity1.Year)
				assert.Equal(t, carDetails, *pinkSlip.Entity1.Details)
			})
			t.Run("LoadEntities", func(t *testing.T) {
				// reset the pink slip
				pinkSlip = &exampleLib.PinkSlip{
					DiLink: *dynamo.NewDiLink(minimalUser, minimalCar),
				}
				userLoaded, carLoaded, err := pinkSlip.LoadEntities(ctx)
				assert.Equal(t, nil, err)
				assert.Equal(t, true, userLoaded)
				assert.Equal(t, true, carLoaded)
				assert.Equal(t, email, pinkSlip.Entity0.Email)
				assert.Equal(t, name, pinkSlip.Entity0.Name)
				assert.Equal(t, age, pinkSlip.Entity0.Age)
				assert.Equal(t, carmake, pinkSlip.Entity1.Make)
				assert.Equal(t, model, pinkSlip.Entity1.Model)
				assert.Equal(t, year, pinkSlip.Entity1.Year)
				assert.Equal(t, carDetails, *pinkSlip.Entity1.Details)
			})
		})
	})
}
