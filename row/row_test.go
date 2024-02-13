package row

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRow(t *testing.T) {
	os.Setenv("TABLE_NAME", "arctica")
	t.Run("Save the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
			Name:  "Test2",
		})
		_, err := userRow.Put(tstCtx, nil)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Save the user, receiving old user as return value", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
			Name:  "Test",
		})
		oldUser, err := userRow.Put(tstCtx, nil)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Test2", oldUser.Object().Name)
	})
	t.Run("Get the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
		})
		err := userRow.Get(tstCtx)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Test", userRow.Object().Name)
		assert.Equal(t, "/r#(User)/rPk(test@gmail.com)", userRow.Keys.Pk)
		assert.Equal(t, "/rSk(details)", userRow.Keys.Sk)
		assert.Equal(t, "/r#(User)/rPk(Test)", *userRow.Keys.Pk1)
		assert.Equal(t, "/rSk(details)", *userRow.Keys.Sk1)
	})
	t.Run("Delete the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
		})
		err := userRow.Delete(tstCtx)
		if err != nil {
			t.Error(err)
		}
		assert.Nil(t, err)
		assert.Equal(t, "Test", userRow.Object().Name)
	})

}
