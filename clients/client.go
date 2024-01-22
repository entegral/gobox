package clients

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

// Client is the primary export of this module
type Client struct {
	Config      aws.Config
	tablename   string
	s3          *awsS3.Client
	dynamo      *dynamodb.Client
	sqs         *sqs.Client
	eventBridge *eventbridge.Client
	http        http.Client
}

// TableName returns the name of the DynamoDB table.
func TableName(ctx context.Context) string {
	tn := os.Getenv("TABLENAME")
	if tn == "" {
		panic("TABLENAME environment variable not set")
	}
	return tn
}

// TableName returns the name of the DynamoDB table.
func (c *Client) TableName(ctx context.Context) string {
	if c.tablename == "" {
		return TableName(ctx)
	}
	return c.tablename
}

var defaultClient *Client

// GetDefaultClient returns the singleton client. If the singleton
// client does not exist, it will be created using the default config
func GetDefaultClient(ctx context.Context) *Client {
	if defaultClient == nil {
		c := newClient(ctx)
		defaultClient = &c
	}
	return defaultClient
}

// SetDefaultClient sets the provided client as the singleton client
func SetDefaultClient(ctx context.Context, client Client) {
	defaultClient = &client
}

// new creates a new client appropriate for use when running within an aws service (i.e. lambda)
func newClient(ctx context.Context) Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	fmt.Println(cfg.Region)
	if err != nil {
		logrus.Fatal("error", err)
	}
	return NewFromConfig(ctx, cfg)
}

// SetDefaultClientToLocalStack sets the default client to a LocalStack client
func SetDefaultClientToLocalStack(ctx context.Context) {
	c := newLocalStackClient(ctx)
	defaultClient = &c
}

// newLocalStackClient creates a new client configured to use LocalStack
func newLocalStackClient(ctx context.Context) Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://host.docker.internal:4566",
				SigningRegion: "us-west-2",
			}, nil
		})),
	)
	if err != nil {
		logrus.Fatal("error", err)
	}
	return NewFromConfig(ctx, cfg)
}

// NewLongTermCredentialClient creates a new client appropriate for use when you have long term credentials
func NewLongTermCredentialClient(ctx context.Context, accessKeyID, secretAccessKey string) Client {
	cfg, err := newConfigWithLongTermCredentials(ctx, accessKeyID, secretAccessKey)
	if err != nil {
		logrus.Fatal("error", err)
	}
	return NewFromConfig(ctx, cfg)
}

// NewClientFromSTSCredentials creates a new client appropriate for use when you have short term STS
// credentials like accessKeyID, secretAccessKey, and sessionToken.
func NewClientFromSTSCredentials(ctx context.Context, accessKeyID, secretAccessKey, sessionToken string) Client {
	cfg, err := newConfigWithCredentials(ctx, accessKeyID, secretAccessKey, sessionToken)
	if err != nil {
		logrus.Fatal("error", err)
	}
	return NewFromConfig(ctx, cfg)
}

// NewFromConfig creates a new client when you have a custom config that you want to use.
func NewFromConfig(ctx context.Context, cfg aws.Config) Client {
	return Client{
		Config: cfg,
		http: http.Client{
			Timeout: 30,
		},
	}
}

func (c Client) WithTableName(tablename string) Client {
	c.tablename = tablename
	return c
}

func (c Client) WithConfig(cfg aws.Config) Client {
	c.Config = cfg
	return c
}

// newConfigWithCredentials creates an AWS Config using the provided IAM credentials.
func newConfigWithCredentials(ctx context.Context, accessKeyID, secretAccessKey, sessionToken string) (aws.Config, error) {
	// Create a static credentials provider
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, sessionToken)
	// Load the default configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

// NewConfigWithLongTermCredentials creates an AWS Config using the provided IAM credentials.
// The session token is optional because it is only required for temporary credentials.
func newConfigWithLongTermCredentials(ctx context.Context, accessKeyID, secretAccessKey string) (aws.Config, error) {
	return newConfigWithCredentials(ctx, accessKeyID, secretAccessKey, "")
}

// S3 returns the s3 client, or creates one if one doesnt exist
func (c *Client) S3() *awsS3.Client {
	if c.s3 == nil {
		c.s3 = s3.NewFromConfig(c.Config)
	}
	return c.s3
}

type DynamoMethods interface {
	GetItem(ctx context.Context, in *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, in *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	UpdateItem(ctx context.Context, in *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

// Dynamo returns the Dynamo client, or creates one if one doesnt exist
func (c *Client) Dynamo() DynamoMethods {
	if c.dynamo == nil {
		c.dynamo = dynamodb.NewFromConfig(c.Config)
	}
	return c.dynamo
}

// SQS returns the SQS client, or creates one if one doesnt exist
func (c *Client) SQS() *sqs.Client {
	if c.sqs == nil {
		c.sqs = sqs.NewFromConfig(c.Config)
	}
	return c.sqs
}

// EventBridge returns the EventBridge client, or creates one if one doesnt exist
func (c *Client) EventBridge() *eventbridge.Client {
	if c.eventBridge == nil {
		c.eventBridge = eventbridge.NewFromConfig(c.Config)
	}
	return c.eventBridge
}
