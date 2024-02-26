package dynamo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/entegral/gobox/clients"
)

type ItemKey struct {
	PkPrefix string
	SkPrefix string
	Index    int
	PkValue  string
	SkValue  string
}

func (i *ItemKey) QueryInput(tableName string) *dynamodb.QueryInput {
	keyConditionExpression := i.KeyConditionExpression()
	expAttrValues := i.ExpressionAttributeValues()
	expAttrKeys := i.ExpressionAttributeKeys()

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expAttrValues,
		ExpressionAttributeNames:  expAttrKeys,
	}

	// Only set the IndexName field if not querying the primary index
	if indexName := i.IndexName(); indexName != "" {
		queryInput.IndexName = aws.String(indexName)
	}

	return queryInput
}

// IndexName returns the name of the Global Secondary Index (GSI) to query.
// It uses the PkPrefix and SkPrefix fields to determine the index name.
// If PkPrefix and SkPrefix are not empty, it will return the index name as: ${PkPrefix}-${SkPrefix}-index
// If PkPrefix and SkPrefix are empty and the Index > 0, it will use the Index field to determine the index name using: "pk${Index}-sk${Index}-index"
// If both PkPrefix and SkPrefix are empty and Index is 0, it will return an empty strini.
func (i *ItemKey) IndexName() string {
	if i.PkPrefix != "" && i.SkPrefix != "" {
		return fmt.Sprintf("%s-%s-index", i.PkPrefix, i.SkPrefix)
	}

	if i.Index > 0 {
		return fmt.Sprintf("pk%d-sk%d-index", i.Index, i.Index)
	}

	return ""
}

func (i *ItemKey) KeyConditionExpression() string {
	return "#pk = :pk AND begins_with(#sk, :sk)"
}

func (i *ItemKey) ExpressionAttributeKeys() map[string]string {
	if i.PkPrefix != "" && i.SkPrefix != "" {
		return map[string]string{
			"#pk": i.PkPrefix,
			"#sk": i.SkPrefix,
		}
	}

	if i.PkPrefix == "" && i.SkPrefix == "" && i.Index > 0 {
		return map[string]string{
			"#pk": "pk" + strconv.Itoa(i.Index),
			"#sk": "sk" + strconv.Itoa(i.Index),
		}
	}

	if i.PkPrefix == "" && i.SkPrefix == "" && i.Index == 0 {
		return map[string]string{
			"#pk": "pk",
			"#sk": "sk",
		}
	}
	return map[string]string{}
}

func (i *ItemKey) ExpressionAttributeValues() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		":pk": &types.AttributeValueMemberS{Value: i.PkValue},
		":sk": &types.AttributeValueMemberS{Value: i.SkValue},
	}
}

type TableQuery struct {
	Input dynamodb.QueryInput
	*dynamodb.QueryOutput
}

type Table struct {
	lock sync.Mutex

	Client    *clients.Client
	Tablename string

	Queries []TableQuery `dynamodbav:"-" json:"-"`
}

// Query executes a query against a specific Global Secondary Index (GSI) or the primary index
// using a provided GSI struct. It checks the GSI struct for values and uses them to determine which
// index to query.
// - if PkPrefix and SkPrefix are not empty, it will query an index with the following name: ${PkPrefix}-${SkPrefix}-index
// - if PkPrefix and SkPrefix are empty, it will use the Index field to determine which primary index to query using: "pk${Index}-sk${Index}-index"
// - If both PkPrefix and SkPrefix are empty and Index is 0, it will query the primary index.
func (t *Table) Query(ctx context.Context, gsi ItemKey) (*dynamodb.QueryOutput, error) {
	tableName := t.TableName(ctx)

	// Use the QueryInput method to generate the query
	queryInput := gsi.QueryInput(tableName)

	c := clients.GetDefaultClient(ctx)
	// Execute the query
	output, err := c.Dynamo().Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	// Add the output to the Queries slice in a concurrent-safe way
	t.lock.Lock()
	t.Queries = append(t.Queries, TableQuery{Input: *queryInput, QueryOutput: output})
	t.lock.Unlock()

	return output, nil
}

// NextQuery will return the next set of results from the last query.
// If there are no more results, it will return nil.
func (t *Table) NextQuery(ctx context.Context) (*dynamodb.QueryOutput, error) {
	if len(t.Queries) == 0 {
		return nil, nil
	}

	// Get the last query
	lastQuery := t.Queries[len(t.Queries)-1]

	// If the last query does not have a LastEvaluatedKey, return nil
	if lastQuery.LastEvaluatedKey == nil {
		return nil, nil
	}

	// Create a new QueryInput using the LastEvaluatedKey from the last query
	queryInput := lastQuery.Input
	queryInput.ExclusiveStartKey = lastQuery.LastEvaluatedKey

	// Execute the query
	c := clients.GetDefaultClient(ctx)
	output, err := c.Dynamo().Query(ctx, &queryInput)
	if err != nil {
		return nil, err
	}

	// Add the output to the Queries slice in a concurrent-safe way
	t.lock.Lock()
	t.Queries = append(t.Queries, TableQuery{Input: queryInput, QueryOutput: output})
	t.lock.Unlock()

	return output, nil
}

func NewTable(tablename string) Table {
	return Table{Tablename: tablename}
}

// TableName returns the name of the DynamoDB table.
// By default, this is the value of the TABLENAME environment variable.
// If you need to override this, implement this method on the parent type.
func (t *Table) TableName(ctx context.Context) string {
	if t.Tablename != "" {
		return t.Tablename
	}
	tn := os.Getenv("TABLENAME")
	if tn == "" {
		panic("TABLENAME environment variable not set")
	}
	return tn
}

func (t *Table) SetTableName(tablename string) {
	t.Tablename = tablename
}

// SetClient sets the client to use for DynamoDB operations.
func (t *Table) SetClient(client *clients.Client) {
	t.Client = client
}
