package dynamo

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiLink(t *testing.T) {
	// init
	os.Setenv("TESTING", "true")
	os.Setenv("TABLENAME", "arctica")
	// t.Setenv("TABLENAME", os.Getenv("TABLENAME"))
	// t.Setenv("TESTING", os.Getenv("TESTING")) // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	// ensure we start with an empty slate
	const email = "testDiLinkEmail@gmail.com"
	const name = "TestDiLinkName"
	const age = 33
	preClearedUser := &User{
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
	var carDetails = CarDetails{
		"doors": float64(4),
		"color": "black",
	}
	preClearedCar := &Car{
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
	pinkSlip := &PinkSlip{
		DiLink: *NewDiLink(preClearedUser, preClearedCar),
	}
	err = pinkSlip.Delete(ctx, pinkSlip)
	if err != nil {
		t.Error(err)
	}

	// generate some minimal entities for testing
	// these only have the minimum required fields to generate the composite keys
	minimalUser := &User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	minimalCar := &Car{
		Make:  carmake,
		Model: model,
		Year:  year,
	}

	// create a second car for testing
	car2Details := CarDetails{
		"doors": float64(2),
		"color": "red",
	}
	car2 := &Car{
		Make:    "TestDiLinkMake2",
		Model:   "TestDiLinkModel2",
		Year:    2023,
		Details: &car2Details,
	}
	err = car2.Delete(ctx, car2)
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
					assert.Contains(t, err.Error(), "item not found:")
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
					assert.Contains(t, err.Error(), "item not found:")
				})
			})
			t.Run("Should not return an error", func(t *testing.T) {
				// create the car now
				err := preClearedCar.Put(ctx, preClearedCar)
				if err != nil {
					t.Error(err)
				}
				t.Run("Should return false when the link isn't in dynamo", func(t *testing.T) {
					linkExists, err := pinkSlip.CheckLink(ctx, pinkSlip, preClearedUser, preClearedCar)
					assert.Nil(t, err)
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
		t.Run("Load Entities of Link", func(t *testing.T) {
			pinkSlip := &PinkSlip{
				DiLink: *NewDiLink(minimalUser, minimalCar),
			}
			err := minimalUser.Put(ctx, minimalUser)
			if err != nil {
				t.Error(err)
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
				pinkSlip = &PinkSlip{
					DiLink: *NewDiLink(minimalUser, minimalCar),
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
		t.Run("FindEntities", func(t *testing.T) {
			t.Run("LoadEntity0s", func(t *testing.T) {
				// start with a fresh pink slip
				pinkSlip := &PinkSlip{
					DiLink: *NewDiLink(minimalUser, minimalCar),
				}
				t.Run("Should return an array of cars when pink slips do exist", func(t *testing.T) {
					// now the pink slip has been created locally with only enough info to generate the composite keys
					entities, err := pinkSlip.LoadEntity0s(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					t.Log("entities:", entities)
					assert.Equal(t, nil, err)
					assert.Equal(t, 1, len(entities))
					assert.NotEqual(t, nil, entities[0])
					assert.Equal(t, preClearedUser.Email, entities[0].Email)
					assert.Equal(t, preClearedUser.Name, entities[0].Name)
					assert.Equal(t, preClearedUser.Age, entities[0].Age)
				})
				t.Run("Should return an empty array when no pink slips exist", func(t *testing.T) {
					// now we have to delete the pink slip from dynamo and re-run the test
					err := pinkSlip.Delete(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					// now we should get an empty array
					entities, err := pinkSlip.LoadEntity0s(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					t.Log("entities:", entities)
					assert.Equal(t, nil, err)
					assert.Equal(t, 0, len(entities))
				})
			})
			t.Run("FindEntity1s", func(t *testing.T) {
				t.Run("Should return an empty array when no pink slips exist", func(t *testing.T) {
					// we should get a similar result here as we did above
					// now we should get an empty array
					entities, err := pinkSlip.LoadEntity1s(ctx, pinkSlip)
					t.Log("entities:", entities)
					assert.IsType(t, &ErrEntityNotFound[*Car]{}, err)
					assert.Equal(t, 0, len(entities))
				})
				t.Run("Should return an array of cars when pink slips do exist", func(t *testing.T) {
					// should get a list of cars back after saving the pink slip again
					err := pinkSlip.Put(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					err = preClearedCar.Put(ctx, preClearedCar)
					if err != nil {
						t.Error(err)
					}
					err = car2.Put(ctx, car2)
					if err != nil {
						t.Error(err)
					}
					entities, err := pinkSlip.LoadEntity1s(ctx, pinkSlip)
					t.Log("entities:", entities)
					assert.Nil(t, err)
					assert.Equal(t, 2, len(entities))
					assert.Equal(t, preClearedCar.Make, entities[1].Make)
					assert.Equal(t, preClearedCar.Model, entities[1].Model)
					assert.Equal(t, preClearedCar.Year, entities[1].Year)
					assert.Equal(t, preClearedCar.Details, entities[1].Details)
				})
				t.Run("Should return an array of cars when multiple pink slips do exist", func(t *testing.T) {

					err := car2.Put(ctx, car2)
					if err != nil {
						t.Error(err)
					}
					pinkSlip2 := &PinkSlip{
						DiLink: *NewDiLink(preClearedUser, car2),
					}
					err = pinkSlip2.Put(ctx, pinkSlip2)
					if err != nil {
						t.Error(err)
					}
					// now we should get an array of cars
					entities, err := pinkSlip.LoadEntity1s(ctx, pinkSlip)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, nil, err)
					assert.Equal(t, 2, len(entities))
					sort.Slice(entities, func(i, j int) bool {
						return entities[i].Year < entities[j].Year
					})
					t.Log("entities:", entities)
					assert.Equal(t, preClearedCar.Make, entities[0].Make)
					assert.Equal(t, preClearedCar.Model, entities[0].Model)
					assert.Equal(t, preClearedCar.Year, entities[0].Year)
					assert.NotEqual(t, nil, *entities[0].Details)
					assert.Equal(t, carDetails, *entities[0].Details)
					assert.Equal(t, car2.Make, entities[1].Make)
					assert.Equal(t, car2.Model, entities[1].Model)
					assert.Equal(t, car2.Year, entities[1].Year)
					assert.Equal(t, car2Details, *entities[1].Details)
				})
			})
		})
		t.Run("FindLinks", func(t *testing.T) {})
	})
}
