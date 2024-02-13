package row

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Query retrieves rows from DynamoDB that match the provided key.
// It runs in a separate goroutine, making it non-blocking and allowing the function to return immediately.
// This means the calling function can continue executing while the query is performed in the background.
// Results and errors are sent through channels as they become available.
// The key is used to create a QueryInput for DynamoDB, specifying the primary key and sort key.
// If successful, it unmarshals the returned items into Rows and sends them to the results channel.
// If an error occurs at any point, it is sent to the errors channel.
func (item *Row[T]) Query(ctx context.Context, key Key) (<-chan Row[T], <-chan error) {
	// Create the channels for results and errors
	results := make(chan Row[T])
	errs := make(chan error)

	go func() {
		defer close(results)
		defer close(errs)

		// Create the key for DynamoDB
		dynamoKey := map[string]awstypes.AttributeValue{
			":pk": &awstypes.AttributeValueMemberS{
				Value: key.Pk,
			},
			":sk": &awstypes.AttributeValueMemberS{
				Value: key.Sk,
			},
		}

		// Create the QueryInput
		queryInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: dynamoKey,
			KeyConditionExpression:    aws.String("pk = :pk and sk = :sk"),
			TableName:                 item.TableName(),
			IndexName:                 key.IndexName(),
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
