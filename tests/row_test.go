package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/stretchr/testify/assert"
)

func TestRow(t *testing.T) {
	// init
	os.Setenv("TABLENAME", "arctica")
	os.Setenv("GOBOX_TESTING", "true") // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	preclear := &exampleLib.User{
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
			type TimeCapsule struct {
				dynamo.Row
				Name string `dynamo:"name"` // Name of the TimeCapsule
			}
			var timeCapsuleGUID string
			t.Run("by default, guids are generated and set to the pk during Put", func(t *testing.T) {})
			t.Run("Other methods rely on the Pk value being set", func(t *testing.T) {
				t.Run("this tool differs than most beause its methods must accept the caller as an argument", func(t *testing.T) {
					t.Run("Put", func(t *testing.T) {
						tc := &TimeCapsule{
							Name: "testName",
						}
						err := tc.Put(ctx, tc) // see, heres that funny business. little weird, but i think you'll find it worth it when
						if err != nil {        // you see the some of the other stuff this tool can do
							t.Error(err)
						}
						timeCapsuleGUID = tc.Pk
						t.Log("GUID written:", tc.Pk)
						assert.Equal(t, true, tc.WasPutSuccessful())
						assert.Equal(t, float64(1), *tc.PutItemOutput.ConsumedCapacity.CapacityUnits)
					})
					t.Run("Get", func(t *testing.T) {
						tc := &TimeCapsule{}
						tc.Pk = timeCapsuleGUID
						loaded, err := tc.Get(ctx, tc) // ಠ_ಠ  it is what it is, not my favorite either
						if err != nil {
							t.Error(err)
						}
						assert.Equal(t, true, loaded)
						assert.Equal(t, true, tc.WasGetSuccessful())
						assert.Equal(t, "testName", tc.Name)
						assert.Equal(t, float64(0.5), *tc.GetItemOutput.ConsumedCapacity.CapacityUnits)
					})
					t.Run("Delete", func(t *testing.T) {
						tc := &TimeCapsule{}
						tc.Pk = timeCapsuleGUID
						err := tc.Delete(ctx, tc) // ಠ_ಠ
						if err != nil {           // but its pretty cool that you already implemented basic CRUD methods
							t.Error(err) // for your new type (and in a way that scales with dynamodb by default). Good job!
						}
						assert.Equal(t, float64(1), *tc.DeleteItemOutput.ConsumedCapacity.CapacityUnits)
					})
				})
			})
			t.Run("the Pk value can be set manually before use with the other methods", func(t *testing.T) {
				t.Run("Get", func(t *testing.T) {})
				t.Run("Delete", func(t *testing.T) {})
			})

		})
		t.Run("heres an example of how you use the methods", func(t *testing.T) {
			t.Run("put", func(t *testing.T) {
				user := &exampleLib.User{
					Email: "testEmail@gmail.com",
					Name:  "TestName",
					Age:   30,
				}
				err = user.Put(ctx, user)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, true, user.WasPutSuccessful())
				assert.Equal(t, float64(1), *user.PutItemOutput.ConsumedCapacity.CapacityUnits)
			})
			t.Run("get", func(t *testing.T) {
				user := &exampleLib.User{
					Email: "testEmail@gmail.com",
					Name:  "TestName",
					Age:   30,
				}
				u := &exampleLib.User{
					Email: user.Email,
				}
				loaded, err := u.Get(ctx, u)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, true, loaded)
				assert.Equal(t, true, u.WasGetSuccessful())
				assert.Equal(t, user.Email, u.Email)
				assert.Equal(t, user.Name, u.Name)
				assert.Equal(t, user.Age, u.Age)
				assert.Equal(t, float64(0.5), *u.GetItemOutput.ConsumedCapacity.CapacityUnits)
			})
			t.Run("delete", func(t *testing.T) {
				user := &exampleLib.User{
					Email: "testEmail@gmail.com",
					Name:  "TestName",
					Age:   30,
				}
				err := user.Put(ctx, user)
				if err != nil {
					t.Error(err)
				}
				if !user.WasPutSuccessful() {
					slog.Error("success:", "WasPutSuccessful", user.WasPutSuccessful())
					slog.Error("got:", "PutItemOutput", user.PutItemOutput)
					t.Error("user put required for delete test")
				}
				u := &exampleLib.User{
					Email: user.Email,
				}
				err = u.Delete(ctx, u)
				if err != nil {
					t.Error(err)
				}
				if len(u.OldDeleteValues()) == 0 {
					slog.Error("success:", "OldDeleteValues", user.OldDeleteValues())
					t.Error("user.OldDeleteValues() should not be empty after delete")
				}
				assert.Equal(t, float64(1), *u.DeleteItemOutput.ConsumedCapacity.CapacityUnits)
			})
		})
	})

	t.Run("it will use a pkshard-sk-index for everything for migrations and to enable querying for backups", func(t *testing.T) {})
	t.Run("TTL", func(t *testing.T) {
		t.Run("it will respect the value configured on the dynamoDBOperations.TTL field", func(t *testing.T) {})
		t.Run("the value will be a pointer to a time.Time instance", func(t *testing.T) {})
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
				dynamo.Row
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
