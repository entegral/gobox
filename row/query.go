package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (item *Row[T]) Query(ctx context.Context, key Key) (<-chan Row[T], <-chan error) {
	// Create the channels for results and errors
	results := make(chan Row[T])
	errs := make(chan error)

	go func() {
		defer close(results)
		defer close(errs)

		// Get the index name from the key
		indexName := GetIndexName(key)

		// Create the key for DynamoDB
		dynamoKey := map[string]awstypes.AttributeValue{
			":pk": &awstypes.AttributeValueMemberS{
				Value: key.PK,
			},
			":sk": &awstypes.AttributeValueMemberS{
				Value: key.SK,
			},
		}

		// Create the QueryInput
		queryInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: dynamoKey,
			KeyConditionExpression:    aws.String("pk = :pk and sk = :sk"),
			TableName:                 aws.String(item.TableName()),
			IndexName:                 aws.String(indexName),
		}

		// Call DynamoDB Query
		result, err := item.GetClient(ctx).Dynamo().Query(ctx, queryInput)
		if err != nil {
			errs <- err
			return
		}

		// Unmarshal the results into Rows and send them to the results channel
		for _, itemMap := range result.Items {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				var row Row[T]
				if err := attributevalue.UnmarshalMap(itemMap, &row); err != nil {
					errs <- err
					return
				}
				results <- row
			}
		}
	}()

	return results, errs
}
