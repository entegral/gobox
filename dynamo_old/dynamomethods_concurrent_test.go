package dynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/entegral/gobox/types"
	"github.com/stretchr/testify/assert"
)

type mockLinkable struct {
	Row
	Name string `json:"name"`
}

func (m *mockLinkable) Keys(gsi int) (partitionKey, sortKey string, err error) {
	return "mockLinkablePartitionKey", "mockLinkableSortKey", nil
}

func (m *mockLinkable) Type() string {
	return "mockLinkable"
}

func TestBatchGet(t *testing.T) {
	db := &DBManager{}
	rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}}

	results := db.BatchGet(context.Background(), rows)

	for result := range results {
		if result.Error != nil {
			assert.IsType(t, &ErrItemNotFound{}, result.Error)
		} else {
			assert.True(t, result.Loaded)
		}
	}
}

func TestBatchPut(t *testing.T) {
	db := &DBManager{}
	rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}}

	results := db.BatchPut(context.Background(), rows)

	for result := range results {
		assert.NoError(t, result.Error)
	}
}

func TestBatchDelete(t *testing.T) {
	db := &DBManager{}
	rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}}

	results := db.BatchDelete(context.Background(), rows)

	for result := range results {
		assert.NoError(t, result.Error)
	}
}

func TestBatchLoadFromMessage(t *testing.T) {
	ml := mockLinkable{
		Name: "loadFromSQSTestName",
	}
	t.Run("should return ErrSQSMessageEmpty if no body is present in message", func(t *testing.T) {
		db := &DBManager{}
		messages := []sqstypes.Message{{}, {}}
		rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}}

		results := db.BatchLoadFromMessage(context.Background(), messages, rows)

		for result := range results {
			if result.Error != nil {
				assert.IsType(t, &ErrSQSMessageEmpty{}, result.Error)
			} else {
				assert.True(t, result.Loaded)
			}
		}
	})
	t.Run("should return Loaded == true and have loaded the item from dynamo", func(t *testing.T) {

		err := ml.Put(context.Background(), &ml)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			err := ml.Delete(context.Background(), &ml)
			if err != nil {
				t.Error(err)
			}
		}()

		db := &DBManager{}
		messages := []sqstypes.Message{
			{
				Body: aws.String(`{"Type": "row", "Pk": "ttlTestGUID1", "Sk": "row"}`),
			},
		}

		rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}, &Row{}}

		results := db.BatchLoadFromMessage(context.Background(), messages, rows)

		for result := range results {
			assert.True(t, result.Loaded)
			assert.NoError(t, result.Error)
			assert.Equal(t, "loadFromSQSTestName", rows[result.Index].(*mockLinkable).Name)
		}
	})
	t.Run("should return ErrItemNotFound if body is present but item fails to load", func(t *testing.T) {
		err := ml.Delete(context.Background(), &ml)
		if err != nil {
			t.Error(err)
		}
		db := &DBManager{}
		messages := []sqstypes.Message{
			{
				Body: aws.String(`{"Type": "mockLinkable", "Pk": "nonexistent", "Sk": "nonexistent"}`),
			},
			{
				Body: aws.String(`{"Type": "mockLinkable", "Pk": "nonexistent", "Sk": "nonexistent"}`),
			},
		}
		rows := []types.Linkable{&mockLinkable{}, &mockLinkable{}, &Row{}}

		results := db.BatchLoadFromMessage(context.Background(), messages, rows)

		for result := range results {
			assert.IsType(t, &ErrItemNotFound{}, result.Error)
		}
	})
}
