package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/entegral/gobox/examples/exampleLib"
	"github.com/stretchr/testify/assert"
)

func TestRow(t *testing.T) {
	// init
	os.Setenv("TABLENAME", "arctica")
	os.Setenv("TESTING", "true") // this will ensure return consumed capacity values are returned
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

	t.Run("CRUD", func(t *testing.T) {
		t.Run("user", func(t *testing.T) {
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
}
