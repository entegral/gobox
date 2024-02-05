package dynamo

import (
	"context"
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
	default:
		return "", "", nil
	}
}

func (u *User) Type() string {
	return "User"
}

var ctx = context.Background()

func TestRow(t *testing.T) {
	t.Run("Save the user", func(t *testing.T) {
		user := User{
			Email: "test@gmail.com",
			Name:  "Test",
		}
		old, err := Put(ctx, &user)
		if err != nil {
			t.Error(err)
		}
		t.Log(old)
	})
	t.Run("Get the user", func(t *testing.T) {
		user := User{
			Email: "test@gmail.com",
		}
		err := Get(ctx, &user)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Test", user.Name)
	})
}
