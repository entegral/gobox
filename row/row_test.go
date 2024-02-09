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

func (u *User) Keys(gsi int) (Key, error) {
	switch gsi {
	case 0:
		return Key{
			PK:    u.Email,
			SK:    "details",
			Index: 0,
		}, nil
	case 1:
		return Key{
			PK:    u.Name,
			SK:    "details",
			Index: 1,
		}, nil
	default:
		return Key{}, nil
	}
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
			Name:  "Test",
		})
		old, err := userRow.Put(ctx, nil)
		if err != nil {
			t.Error(err)
		}
		t.Log(old)
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
		assert.Equal(t, userRow.object.Email, userRow.Keys.PrimaryKey.PK)
		assert.Equal(t, "details", userRow.Keys.PrimaryKey.SK)
		assert.Equal(t, "Test", userRow.Keys.GSI1.PK)
		assert.Equal(t, "details", userRow.Keys.GSI1.SK)
	})
	t.Run("Delete the user", func(t *testing.T) {
		userRow := NewRow(&User{
			Email: "test@gmail.com",
		})
		err := userRow.Delete(ctx)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Test", userRow.Object().Name)
	})

}
