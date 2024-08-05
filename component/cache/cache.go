package cache

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBClient(sess *session.Session, tableName string) (*DynamoDBClient, error) {
	client := dynamodb.New(sess)

	return &DynamoDBClient{
		client:    client,
		tableName: tableName,
	}, nil
}

func (d *DynamoDBClient) Set(key string, value string) error {
	partitionKey := os.Getenv("DYNAMODB_PARTITION_KEY")
	item := map[string]*dynamodb.AttributeValue{
		partitionKey: {
			S: aws.String(key),
		},
		"value": {
			S: aws.String(value),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.tableName),
	}

	_, err := d.client.PutItem(input)
	return err
}

func (d *DynamoDBClient) Get(key string) (string, error) {
	partitionKey := os.Getenv("DYNAMODB_PARTITION_KEY")
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(key),
			},
		},
		TableName:      aws.String(d.tableName),
		ConsistentRead: aws.Bool(true),
	}

	result, err := d.client.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", nil
	}

	value, ok := result.Item["value"]
	if !ok || value.S == nil {
		return "", nil
	}

	return *value.S, nil
}
