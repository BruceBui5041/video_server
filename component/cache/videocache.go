package cache

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"video_server/appconst"
	"video_server/logger"
	models "video_server/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
)

func (d *DynamoDBClient) SetVideoCache(courseSlug string, video models.Video) error {
	videoTableName := os.Getenv("DYNAMODB_VIDEO_TABLE_NAME")
	cacheKey := fmt.Sprintf("%s:%s:%d", appconst.VideoURLPrefix, courseSlug, video.Id)
	item := map[string]*dynamodb.AttributeValue{
		"cachekey": {
			S: aws.String(cacheKey),
		},
		"videoid": {
			N: aws.String(strconv.FormatInt(int64(video.Id), 10)),
		},
		"value": {
			S: aws.String(video.VideoURL),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(videoTableName),
	}

	_, err := d.client.PutItem(input)
	return err
}

func (d *DynamoDBClient) GetVideoCache(courseSlug string, videoId uint32) (string, error) {
	videoTableName := os.Getenv("DYNAMODB_VIDEO_TABLE_NAME")
	cacheKey := fmt.Sprintf("%s:%s:%d", appconst.VideoURLPrefix, courseSlug, videoId)
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info("dynamoDB Get duration",
			zap.String("key", cacheKey),
			zap.Float64("duration_ms", durationMs))
	}()

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(cacheKey),
			},
			"videoid": {
				N: aws.String(strconv.FormatUint(uint64(videoId), 10)),
			},
		},
		TableName:      aws.String(videoTableName),
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
		logger.AppLogger.Warn("dynamoDB not found", zap.String("key", cacheKey))
		return "", nil
	}

	return *value.S, nil
}
