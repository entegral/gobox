package skeleton

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Generator is a function type that generates at least one key for a given item.
type Generator func() (keyName, value string, err error)

// Keys generates keys for the item.
func Keys(generators ...Generator) (map[string]string, error) {
	keys := make(map[string]string)
	for _, generator := range generators {
		keyName, value, err := generator()
		if err != nil {
			return nil, err
		}
		keys[keyName] = value
	}
	return keys, nil
}

// DynamoKeyMapV1 generates keys for the item.
func DynamoKeyMapV1(generators ...Generator) (map[string]*dynamodb.AttributeValue, error) {
	keys := make(map[string]*dynamodb.AttributeValue)
	for _, generator := range generators {
		keyName, value, err := generator()
		if err != nil {
			return nil, err
		}
		keys[":"+keyName] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
	return keys, nil
}
