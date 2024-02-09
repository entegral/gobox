package row

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
		old, err := userRow.Put(ctx, func(pii *dynamodb.PutItemInput) {
			pii.ConditionExpression = aws.String("attribute_not_exists(pk)")
		})
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
	})
}
