To standardize the way to use the query lookup method and support 6 GSIs in addition to queries on the primary table's pk and sk, we can create a new helper function that leverages the existing interfaces and types. This helper function will abstract away the details of constructing the query input and will allow developers to easily query the table based on a partition key and an optional sort key condition.

We will create a new file named `queryhelper.go` in the `helpers` package with the following content:

```go
package dynamo

import (
        "context"
        "fmt"
        "os"

        "github.com/aws/aws-sdk-go-v2/aws"
        "github.com/aws/aws-sdk-go-v2/service/dynamodb"
        "github.com/entegral/gobox/clients"
        "github.com/entegral/gobox/types"
)

// QueryByGSI executes a query against a specific Global Secondary Index (GSI) or the primary index.
// It uses the default client. If you need to use a specific client, use QueryByGSIWithClient instead.
func QueryByGSI[T types.Keyable](ctx context.Context, gsi int, partitionKey, sortKeyCondition string) ([]T, error) {
        client := clients.GetDefaultClient(ctx)
        return QueryByGSIWithClient[T](ctx, client, gsi, partitionKey, sortKeyCondition)
}

// QueryByGSIWithClient executes a query against a specific Global Secondary Index (GSI) or the primary index using the provided client.
func QueryByGSIWithClient[T types.Keyable](ctx context.Context, client *clients.Client, gsi int, partitionKey, sortKeyCondition string) ([]T, error) {
        tableName := os.Getenv("TABLENAME")
        if tableName == "" {
                return nil, fmt.Errorf("TABLENAME environment variable not set")
        }

        var indexName *string
        var keyConditionExpression string
        var expressionAttributeNames map[string]string
        var expressionAttributeValues map[string]types.AttributeValue

        // Construct the key condition expression based on the GSI
        if gsi == 0 {
                keyConditionExpression = "pk = :pk"
                if sortKeyCondition != "" {
                        keyConditionExpression += " AND sk = :sk"
                }
        } else {
                indexName = aws.String(fmt.Sprintf("pk%d-sk%d-index", gsi, gsi))
                keyConditionExpression = fmt.Sprintf("pk%d = :pk", gsi)
                if sortKeyCondition != "" {
                        keyConditionExpression += fmt.Sprintf(" AND sk%d = :sk", gsi)
                }
                expressionAttributeNames = map[string]string{
                        fmt.Sprintf("#pk%d", gsi): fmt.Sprintf("pk%d", gsi),
                        fmt.Sprintf("#sk%d", gsi): fmt.Sprintf("sk%d", gsi),
                }
        }

        // Construct the expression attribute values
        expressionAttributeValues = map[string]types.AttributeValue{
                ":pk": &types.AttributeValueMemberS{Value: partitionKey},
        }
        if sortKeyCondition != "" {
                expressionAttributeValues[":sk"] = &types.AttributeValueMemberS{Value: sortKeyCondition}
        }

        // Create the query input
        queryInput := &dynamodb.QueryInput{
                TableName:                 aws.String(tableName),
                IndexName:                 indexName,
                KeyConditionExpression:    aws.String(keyConditionExpression),
                ExpressionAttributeNames:  expressionAttributeNames,
                ExpressionAttributeValues: expressionAttributeValues,
        }

        // Execute the query
        return QueryListOfMaps[T](ctx, client, *queryInput)
}
```

This new helper function, `QueryByGSI`, takes a generic type `T` that must implement the `types.Keyable` interface. It accepts a context, a GSI number (0 for the primary index), a partition key, and an optional sort key condition. It constructs the necessary query input and delegates to the existing `QueryListOfMaps` function to execute the query and unmarshal the results.

By using this helper, developers can easily query any of the supported GSIs or the primary index with minimal implementation work. They just need to ensure that their types implement the `types.Keyable` interface and then call `QueryByGSI` with the appropriate parameters.