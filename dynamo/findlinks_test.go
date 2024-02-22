package dynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Dynamo() *dynamodb.Client {
	return &dynamodb.Client{}
}

func (m *MockClient) Query(ctx context.Context, params *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func TestFindLinkRowsByEntityGSI(t *testing.T) {
	ctx := context.Background()
	clients := clients.GetDefaultClient(ctx)
	// create user
	user := CreateUser("FindLinkRowsTestUser@gmail.com")
	err := user.Put(ctx, user)
	assert.NoError(t, err)
	defer user.Delete(ctx, user)
	// create car
	car := PutCar(ctx, "FindLinkRowsTestCar", "FindLinkRowsTestModel")
	err = car.Put(ctx, car)
	assert.NoError(t, err)
	defer car.Delete(ctx, car)
	// create link
	pinkSlip := &PinkSlip{
		DiLink: *NewDiLink(user, car),
		VIN:    "123456",
	}
	err = pinkSlip.Put(ctx, pinkSlip)
	assert.NoError(t, err)
	defer pinkSlip.Delete(ctx, pinkSlip)

	t.Run("test we can find the pinkSlip by an Entity0 lookup", func(t *testing.T) {
		items, err := findLinkRowsByEntityGSI(ctx, clients, user, Entity0GSI, pinkSlip.Type())

		expectedPk := "/rowType(PinkSlip)/rowPk(/e0Type(user)/e0pk(FindLinkRowsTestUser@gmail.com)/e1Type(car)/e1pk(FindLinkRowsTestCar-FindLinkRowsTestModel))"
		expectedSk := "/e0sk(info)/e1sk(2018)"
		assert.NoError(t, err)
		assert.Equal(t, 1, len(items))
		assert.Equal(t, expectedPk, items[0]["pk"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, expectedSk, items[0]["sk"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, pinkSlip.Type(), items[0]["type"].(*types.AttributeValueMemberS).Value)
	})
	t.Run("should find the pinkSlip by an Entity1 lookup", func(t *testing.T) {
		items, err := findLinkRowsByEntityGSI(ctx, clients, car, Entity1GSI, pinkSlip.Type())

		expectedPk := "/rowType(PinkSlip)/rowPk(/e0Type(user)/e0pk(FindLinkRowsTestUser@gmail.com)/e1Type(car)/e1pk(FindLinkRowsTestCar-FindLinkRowsTestModel))"
		expectedSk := "/e0sk(info)/e1sk(2018)"
		assert.NoError(t, err)
		assert.Equal(t, 1, len(items))
		assert.Equal(t, expectedPk, items[0]["pk"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, expectedSk, items[0]["sk"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, pinkSlip.Type(), items[0]["type"].(*types.AttributeValueMemberS).Value)
	})
}
