package cache

import (
	"os"
	"time"
	"video_server/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
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
	logger.AppLogger.Debug("dynamoDB set", zap.Any("key", key), zap.Any("value", value))
	return err
}

func (d *DynamoDBClient) Get(key string) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info("dynamoDB Get duration",
			zap.String("key", key),
			zap.Float64("duration_ms", durationMs))
	}()

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
		logger.AppLogger.Debug("dynamoDB not found", zap.String("key", key))
		return "", nil
	}
	logger.AppLogger.Debug("dynamoDB found", zap.String("key", key), zap.String("value", *value.S))
	return *value.S, nil
}
