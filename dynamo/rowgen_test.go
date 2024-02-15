// FILEPATH: /home/robertbruce/repos/gobox/dynamo/rowgen_test.go

package dynamo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserNoEmbed struct {
	Shard
	Table

	Email string
	Name  string
	Age   int
}

func (u *UserNoEmbed) Type() string {
	return "UserNoEmbed"
}

// Keys returns the partition key and sort key for the row
func (u *UserNoEmbed) Keys(gsi int) (string, string, error) {
	if u == nil {
		return "", "", errors.New("nil user")
	}
	if u.Email == "" {
		return "", "", &ErrMissingEmail{}
	}
	switch gsi {
	default:
		// Handle other GSIs or return an error
		return u.Email, "info", nil
	}
}

var test1 = UserNoEmbed{
	Email: "test@example.com",
	Name:  "Test User",
	Age:   30,
}

var test2 = UserNoEmbed{
	Email: "test@example.com",
	Name:  "New Test User",
	Age:   31,
}

var unloadedUser = UserNoEmbed{
	Email: "test@example.com",
}

func TestNewRow(t *testing.T) {
	ctx := context.Background()
	t.Run("simple put and get test", func(t *testing.T) {
		rowGen := NewRow(&test1)
		if rowGen.Type() != "UserNoEmbed" {
			t.Errorf("NewRow().object.Type() = %v, want %v", rowGen.Type(), "UserNoEmbed")
		}
		_, err := rowGen.Put(ctx)
		if err != nil {
			t.Errorf("NewRow().Put() = %v, want %v", err, nil)
		}
		start := unloadedUser
		getRow := NewRow(&start)
		loaded, err := getRow.Get(ctx)
		if err != nil {
			t.Errorf("NewRow().Get() = %v, want %v", err, nil)
		}
		if !loaded {
			t.Errorf("NewRow().Get() = %v, want %v", loaded, true)
		}
		assert.Equal(t, getRow.Object.Age, 30)
	})
	t.Run("test overwrites and delete", func(t *testing.T) {
		t.Run("test overwrites", func(t *testing.T) {
			rowGen := NewRow(&test2)
			oldItem, err := rowGen.Put(ctx)
			if err != nil {
				t.Errorf("NewRow().Put() = %v, want %v", err, nil)
			}
			assert.Equal(t, oldItem.Age, 30)
			assert.Equal(t, oldItem.Name, "Test User")
		})
		t.Run("test delete", func(t *testing.T) {
			start := unloadedUser
			rowGen := NewRow(&start)
			oldItem, err := rowGen.Delete(ctx)
			if err != nil {
				t.Errorf("NewRow().Delete() = %v, want %v", err, nil)
			}
			assert.Equal(t, oldItem.Age, 31)
			assert.Equal(t, oldItem.Name, "New Test User")
		})
	})

	t.Run("test WithTable", func(t *testing.T) {
		rowGen := NewRow(&test1)
		defer rowGen.Delete(ctx)
		old, err := rowGen.WithTableName("test").Put(ctx)
		assert.Error(t, err)
		assert.Nil(t, old)
		old, err = rowGen.WithTableName("arctica").Put(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, old.Age)
		rowGen.Object.Age = 32
		old, err = rowGen.Put(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 30, old.Age)
		assert.Equal(t, 32, rowGen.Object.Age)
	})
}
