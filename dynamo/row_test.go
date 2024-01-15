package dynamo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestRow(t *testing.T) {
	// init

	// os.Setenv("TABLENAME", "arctica")
	// os.Setenv("GOBOX_TESTING", "true") // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	preclear := &User{
		Email: "testEmail@gmail.com",
		Name:  "TestName",
		Age:   30,
	}
	// ensure we start with an empty slate
	err := preclear.Delete(ctx, preclear)
	if err != nil {
		t.Error(err)
	}

	t.Run("CRUD is the bread and butter of the Row type", func(t *testing.T) {
		t.Run("out of the gate, embedding the Row into your type grants you free Get, Put and Update methods", func(t *testing.T) {
			// lets test this with a TimeCapsule struct
			var timeCapsuleGUID string
			t.Run("CRUD methods used in this repo rely on the Pk value being set", func(t *testing.T) {
				type TimeCapsule struct {
					Row
					Name string `dynamo:"name"` // Name of the TimeCapsule
				}
				t.Run("this tool differs than most beause its methods must accept the caller as an argument", func(t *testing.T) {
					t.Run("Put", func(t *testing.T) {
						tc := &TimeCapsule{
							Name: "testName",
						}
						err := tc.Put(ctx, tc) // ಠ_ಠ  see, heres that funny business. little weird, but i think you'll find it worth it when
						if err != nil {        // you see the some of the other stuff this tool can do
							t.Error(err)
						}
						// by default, a guid was generated and set to the pk during the Put call
						// this is the default behavior, but you can override it if you satisfy
						// the Keyable interface, which we'll get to later.
						timeCapsuleGUID = tc.Pk
						t.Log("GUID written:", tc.Pk)
						assert.Equal(t, true, tc.WasPutSuccessful())
						assert.Equal(t, float64(2), *tc.PutItemOutput.ConsumedCapacity.CapacityUnits)
					})
				})
				t.Run("the Pk value can be set manually before use with the other methods", func(t *testing.T) {
					t.Run("Get", func(t *testing.T) {
						tc := &TimeCapsule{}
						// by default, the only field you need to assign is the Pk field to the object's guid
						// after this get and delete call we'll go over how to override this behavior
						tc.Pk = timeCapsuleGUID
						loaded, err := tc.Get(ctx, tc) // ಠ_ಠ  it is what it is, not my favorite either
						if err != nil {
							t.Error(err)
						}
						t.Log("GUID written:", tc.Pk)
						assert.Equal(t, true, loaded)
						assert.Equal(t, true, tc.WasGetSuccessful())
						assert.Equal(t, "testName", tc.Name)
						assert.Equal(t, float64(0.5), *tc.GetItemOutput.ConsumedCapacity.CapacityUnits)
					})
					t.Run("Delete", func(t *testing.T) {
						tc := &TimeCapsule{}
						tc.Pk = timeCapsuleGUID
						err := tc.Delete(ctx, tc)
						// ಠ_ಠ
						// weird, but its pretty cool that you already implemented basic CRUD methods
						// for your new type (and in a way that scales with dynamodb by default). Good job!
						if err != nil {
							t.Error(err)
						}
						t.Log("GUID written:", tc.Pk)
						assert.Equal(t, float64(2), *tc.DeleteItemOutput.ConsumedCapacity.CapacityUnits)
					})
				})
			})

			t.Run("Now we will implement our object's own Keyable interface", func(t *testing.T) {
				// this allows us to override the default GUID behavior of the Row type
				ktc := &KeyableTimeCapsule{
					Name:     "testName",
					Location: "testLocation",
				}
				t.Run("the Keyable interface requires a Keys method", func(t *testing.T) {
					t.Run("the Keys method must return a pk and sk value", func(t *testing.T) {
						pk, sk, err := ktc.Keys(0)
						if err != nil {
							t.Error(err)
						}
						assert.Equal(t, "testLocation", pk)
						assert.Equal(t, "testName", sk)
						t.Run("We must also confirm that the Keyable implementation sets the pk and sk values on the object", func(t *testing.T) {
							assert.Equal(t, "testLocation", ktc.Pk)
							assert.Equal(t, "testName", ktc.Sk)
						})
					})
				})
			})
		})
		t.Run("heres another few examples of how you use the methods", func(t *testing.T) {
			testUserEmail := "testEmail@gmail.com"
			testUserName := "TestName"
			testUserAge := 30
			t.Run("put", func(t *testing.T) {
				// the example User type implements a Keyable function that will use the Email field as the pk
				// and a static "info" string as the sk. Using a constructor function is a good way to ensure
				// that the fields required for the composite key are set.
				user := CreateUser(testUserEmail)
				user.Name = testUserName
				user.Age = testUserAge

				err = user.Put(ctx, user)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, true, user.WasPutSuccessful())
				assert.Equal(t, float64(2), *user.PutItemOutput.ConsumedCapacity.CapacityUnits)
			})
			t.Run("get", func(t *testing.T) {
				u := CreateUser(testUserEmail)

				loaded, err := u.Get(ctx, u)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, true, loaded)
				assert.Equal(t, true, u.WasGetSuccessful())
				assert.Equal(t, testUserEmail, u.Email)
				assert.Equal(t, testUserName, u.Name)
				assert.Equal(t, testUserAge, u.Age)
				assert.Equal(t, float64(0.5), *u.GetItemOutput.ConsumedCapacity.CapacityUnits)
			})
			t.Run("delete", func(t *testing.T) {
				u := CreateUser(testUserEmail)
				err = u.Delete(ctx, u)
				if err != nil {
					t.Error(err)
				}
				if len(u.OldDeleteValues()) == 0 {
					t.Error("user.OldDeleteValues() should not be empty after delete")
				}
				assert.Equal(t, float64(2), *u.DeleteItemOutput.ConsumedCapacity.CapacityUnits)
			})
		})
	})

	t.Run("it will use a pkshard-sk-index for migrations and to enable querying for backups", func(t *testing.T) {
		t.Run("ensure the pkshard attribute is saved to dynamo during a Put", func(t *testing.T) {
			err := preclear.Put(ctx, preclear)
			if err != nil {
				t.Error(err)
			}
			_, err = preclear.Get(ctx, preclear)
			if err != nil {
				t.Error(err)
			}
			pkshardAV := preclear.RowData["pkshard"]
			pkshard, ok := pkshardAV.(*types.AttributeValueMemberS)
			if !ok {
				t.Error("pkshard should be a string")
			}
			assert.Regexp(t, `\.(\d{1,2})$`, pkshard.Value)
		})
	})
	t.Run("TTL", func(t *testing.T) {
		expectedTTL := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		var expectedDynamoTTL int64 = 1577836800
		t.Run("1/1/2020 to unix", func(t *testing.T) {
			assert.Equal(t, expectedDynamoTTL, expectedTTL.Unix())
		})
		t.Run("the DBManager.TTL field will marshal into an int", func(t *testing.T) {
			testGUID := "ttlTestGUID1"
			// lets test this with a TimeCapsule struct
			type TimeCapsule struct {
				Row
				Name string `dynamo:"name"` // Name of the TimeCapsule
			}
			tc := &TimeCapsule{
				Name: "testName",
			}
			tc.Pk = testGUID
			tc.SetTTL(expectedTTL)
			err := tc.Put(ctx, tc)
			if err != nil {
				t.Error(err)
			}
			loadTc := &TimeCapsule{}
			loadTc.Pk = testGUID
			_, err = loadTc.Get(ctx, loadTc)
			if err != nil {
				t.Error(err)
			}
			ttlAV := loadTc.RowData["ttl"]
			ttl, ok := ttlAV.(*types.AttributeValueMemberN)
			if !ok {
				t.Error("ttl should be a number")
			}
			assert.Equal(t, fmt.Sprintf("%d", expectedDynamoTTL), ttl.Value)
			assert.Equal(t, expectedTTL.Unix(), loadTc.TTL.Unix())
		})
		t.Run("the value will not be marshalled to dynamo if unset", func(t *testing.T) {
			testGUID := "ttlTestGUID2"
			// lets test this with a TimeCapsule struct
			type TimeCapsule struct {
				Row
				Name string `dynamo:"name"` // Name of the TimeCapsule
			}
			tc := &TimeCapsule{
				Name: "testName",
			}
			tc.Pk = testGUID
			err := tc.Put(ctx, tc)
			if err != nil {
				t.Error(err)
			}
			loadTc := &TimeCapsule{}
			loadTc.Pk = testGUID
			_, err = loadTc.Get(ctx, loadTc)
			if err != nil {
				t.Error(err)
			}
			assert.Nil(t, loadTc.TTL)
			// tc.Delete(ctx, tc)
		})
		t.Run("the field will only be marshalled to json if the field is set", func(t *testing.T) {})
		t.Run("the field will only be marshalled to a dynamo map if the field is set", func(t *testing.T) {})
	})

	t.Run("Hashed Entropy", func(t *testing.T) {
		/*
			To save types to dynamo, the Keyable interface is used to generate the dynamo composite key.
			By default, this will be done by automatically generating a GUID and assigning it to the pk field
			during a Put call. This will make it impossible to retrieve the row by the pk field without knowing
			the GUID, but for some circumstances this may not matter.

			This will ensure that there is pk entropy by default, however, if you want to override this behavior,
			you can implement the Keyable interface yourself.
		*/
		t.Run("DynamoRow will implement the keyable interface", func(t *testing.T) {
			// lets test this with a TimeCapsule struct
			type TimeCapsule struct {
				Row
				Name string `dynamo:"name"` // Name of the TimeCapsule
			}

			t.Run("when Keyable isnt implemented", func(t *testing.T) {
				testPk := "testPk"
				testSk := "testSk"
				t.Run("it will return the value of Pk and Sk if they are already set", func(t *testing.T) {
					tc := &TimeCapsule{}
					tc.Pk = testPk
					tc.Sk = testSk
					pk, sk, err := tc.Keys(0)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, testPk, pk)
					assert.Equal(t, testSk, sk)
				})
				t.Run("it will generate a GUID for the Pk when unset", func(t *testing.T) {
					tc := &TimeCapsule{}
					pk, _, err := tc.Keys(0)
					if err != nil {
						t.Error(err)
					}
					assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, pk)
				})
				t.Run("it will return the Pk if set, and return a static 'row' value for the Sk when unset", func(t *testing.T) {
					tc := &TimeCapsule{}
					tc.Pk = testPk
					_, sk, err := tc.Keys(0)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, testPk, tc.Pk)
					assert.Equal(t, "row", sk)
				})
				t.Run("it will generate a static 'row' value for the Sk when unset", func(t *testing.T) {
					tc := &TimeCapsule{}
					_, sk, err := tc.Keys(0)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, "row", sk)
				})
			})
		})
	})
}
