package cache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBClient struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBClient(tableName string) (*DynamoDBClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("your-aws-region"),
	})
	if err != nil {
		return nil, err
	}

	client := dynamodb.New(sess)

	return &DynamoDBClient{
		client:    client,
		tableName: tableName,
	}, nil
}

func (d *DynamoDBClient) Set(key string, value string) error {
	item := map[string]*dynamodb.AttributeValue{
		"Key": {
			S: aws.String(key),
		},
		"Value": {
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
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(d.tableName),
	}

	result, err := d.client.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", nil
	}

	var item struct {
		Value string `json:"Value"`
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}

	return item.Value, nil
}
