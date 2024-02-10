package row

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Email string `dynamodbav:"email" json:"email"`
	Name  string `dynamodbav:"name" json:"name"`
}

func (u *User) Keys(gsi int) (string, string, error) {
	switch gsi {
	case 0:
		return u.Email, "details", nil
	case 1:
		return u.Name, "details", nil
	}
	return "", "", nil
}

func (u *User) Type() string {
	return "User"
}

var ctx = context.Background()

func TestRow(t *testing.T) {
	os.Setenv("TABLE_NAME", "arctica")
	t.Run("Save the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
			Name:  "Test2",
		})
		_, err := userRow.Put(ctx, nil)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Save the user, receiving old user as return value", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
			Name:  "Test",
		})
		oldUser, err := userRow.Put(ctx, nil)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Test2", oldUser.Object().Name)
	})
	t.Run("Get the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
		})
		err := userRow.Get(ctx)
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
		err := userRow.Delete(ctx)
		if err != nil {
			t.Error(err)
		}
		assert.Nil(t, err)
		assert.Equal(t, "Test", userRow.Object().Name)
	})

}
