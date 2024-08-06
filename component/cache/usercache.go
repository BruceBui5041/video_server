package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"video_server/appconst"
	"video_server/logger"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func (d *DynamoDBClient) SetUserCache(user models.User) error {
	videoTableName := os.Getenv("DYNAMODB_USER_TABLE_NAME")
	cacheKey := fmt.Sprintf("%s:%d", appconst.UserPrefix, user.Id)

	var cacheUser usermodel.CacheUser
	err := copier.Copy(&cacheUser, user)

	if err != nil {
		return err
	}

	cacheUserJson, err := json.Marshal(cacheUser)

	if err != nil {
		logger.AppLogger.Warn("cannot parse user for cache", zap.Error(err))
		return err
	}

	item := map[string]*dynamodb.AttributeValue{
		"cachekey": {
			S: aws.String(cacheKey),
		},
		"userid": {
			N: aws.String(strconv.FormatInt(int64(user.Id), 10)),
		},
		"value": {
			S: aws.String(string(cacheUserJson)),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(videoTableName),
	}

	_, err = d.client.PutItem(input)
	return err
}

func (d *DynamoDBClient) GetUserCache(userId uint32) (string, error) {
	start := time.Now()
	key := fmt.Sprintf("%s:%d", appconst.UserPrefix, userId)
	userTableName := os.Getenv("DYNAMODB_USER_TABLE_NAME")
	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info("dynamoDB Get duration",
			zap.String("key", key),
			zap.Float64("duration_ms", durationMs))
	}()

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(key),
			},
			"userid": {
				N: aws.String(strconv.FormatUint(uint64(userId), 10)),
			},
		},
		TableName:      aws.String(userTableName),
		ConsistentRead: aws.Bool(false),
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
		logger.AppLogger.Warn("dynamoDB not found", zap.String("key", key))
		return "", nil
	}

	return *value.S, nil
}
