package dynamo

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/clients"
	"github.com/entegral/gobox/types"
)

type Keys struct {
	Pk  string `dynamodbav:"pk,omitempty" json:"pk,omitempty"`
	Sk  string `dynamodbav:"sk,omitempty" json:"sk,omitempty"`
	Pk1 string `dynamodbav:"pk1,omitempty" json:"pk1,omitempty"`
	Sk1 string `dynamodbav:"sk1,omitempty" json:"sk1,omitempty"`
	Pk2 string `dynamodbav:"pk2,omitempty" json:"pk2,omitempty"`
	Sk2 string `dynamodbav:"sk2,omitempty" json:"sk2,omitempty"`
	Pk3 string `dynamodbav:"pk3,omitempty" json:"pk3,omitempty"`
	Sk3 string `dynamodbav:"sk3,omitempty" json:"sk3,omitempty"`
	Pk4 string `dynamodbav:"pk4,omitempty" json:"pk4,omitempty"`
	Sk4 string `dynamodbav:"sk4,omitempty" json:"sk4,omitempty"`
	Pk5 string `dynamodbav:"pk5,omitempty" json:"pk5,omitempty"`
	Sk5 string `dynamodbav:"sk5,omitempty" json:"sk5,omitempty"`
	Pk6 string `dynamodbav:"pk6,omitempty" json:"pk6,omitempty"`
	Sk6 string `dynamodbav:"sk6,omitempty" json:"sk6,omitempty"`
}

func (k *Keys) MaxGSIs() int {
	return 6
}

type Rowable interface {
	types.Keyable
	types.Typeable
}

type row[T Rowable] struct {
	Keys
	object T

	tableName string
	client    *clients.Client
}

func (r *row[T]) UnmarshalMap(m map[string]awstypes.AttributeValue) error {
	// Create a new map to hold the non-key values
	err := attributevalue.UnmarshalMap(m, &r.Keys)
	if err != nil {
		return err
	}

	// Unmarshal the non-key values into the Object
	err = attributevalue.UnmarshalMap(m, &r.object)
	return err
}

func newRow[T Rowable](object T) row[T] {
	return row[T]{object: object}
}

func (r *row[T]) TableName() string {
	// If the table name is not set, use the value provided from TABLE_NAME
	if r.tableName == "" {
		return os.Getenv("TABLENAME")
	}
	return r.tableName
}

func (r *row[T]) SetTableName(tableName string) {
	r.tableName = tableName
}

func (r *row[T]) SetClient(client *clients.Client) {
	r.client = client
}

func (r *row[T]) GetClient(ctx context.Context) *clients.Client {
	// If the client is not set, use the default client from the clients package
	if r.client != nil {
		return r.client
	}
	return clients.GetDefaultClient(ctx)
}
