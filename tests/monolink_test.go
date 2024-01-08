package tests

import (
	"context"
	"os"
	"testing"

	"github.com/entegral/gobox/dynamo"
	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/stretchr/testify/assert"
)

func TestMonoLink(t *testing.T) {
	// init
	os.Setenv("TABLENAME", "arctica")
	os.Setenv("TESTING", "true") // this will ensure return consumed capacity values are returned
	ctx := context.Background()
	email := "testEmail@gmail.com"
	name := "TestName"
	age := 30
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

	const phone = "123-456-7890"
	const addr = "123 Main St"

	t.Run("MonoLink", func(t *testing.T) {
		t.Run("put", func(t *testing.T) {
			// put a user for the test
			user := &exampleLib.User{
				Email: email,
				Name:  name,
				Age:   age,
			}
			err = user.Put(ctx, user)
			if err != nil {
				t.Error(err)
			}
			contact := &exampleLib.ContactInfo{
				MonoLink: dynamo.NewMonoLink(user),
				Phone:    phone,
				Addr:     addr,
			}
			err = contact.Put(ctx, contact)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, true, contact.WasPutSuccessful())
			assert.Equal(t, float64(2), *contact.PutItemOutput.ConsumedCapacity.CapacityUnits)
		})
		t.Run("get link from base", func(t *testing.T) {
			// put a user for the test
			user := &exampleLib.User{
				Email: email,
				Name:  name,
				Age:   age,
			}
			contact := &exampleLib.ContactInfo{
				MonoLink: dynamo.NewMonoLink(user),
			}
			loaded, err := contact.Get(ctx, contact)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, true, loaded)
			assert.Equal(t, true, contact.WasGetSuccessful())
			assert.Equal(t, phone, contact.Phone)
			assert.Equal(t, addr, contact.Addr)
		})
		t.Run("delete link", func(t *testing.T) {
			// put a user for the test
			user := &exampleLib.User{
				Email: email,
				Name:  name,
				Age:   age,
			}
			contact := &exampleLib.ContactInfo{
				MonoLink: dynamo.NewMonoLink(user),
			}
			err := contact.Delete(ctx, contact)
			if err != nil {
				t.Error(err)
			}
			assert.NotEqual(t, nil, contact.DeleteItemOutput)
			assert.Equal(t, float64(2), *contact.DeleteItemOutput.ConsumedCapacity.CapacityUnits)
		})
	})
}
