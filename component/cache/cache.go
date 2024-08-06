package cache

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	client *dynamodb.DynamoDB
}

func NewDynamoDBClient(sess *session.Session) (*DynamoDBClient, error) {
	client := dynamodb.New(sess)

	return &DynamoDBClient{
		client: client,
	}, nil
}
