package storagehandler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
	"video_server/appconst"
	"video_server/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

func GetS3File(svc *s3.S3, bucket, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		logger.AppLogger.Error("Failed to get object from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	logger.AppLogger.Info("Successfully retrieved object from S3", zap.String("bucket", bucket), zap.String("key", key))
	return result.Body, nil
}

func UploadFileToS3(svc *s3.S3, file io.Reader, bucket, key string) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		logger.AppLogger.Error("Failed to read file into buffer", zap.Error(err))
		return fmt.Errorf("failed to read file: %v", err)
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})

	if err != nil {
		logger.AppLogger.Error("Failed to upload file to S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return fmt.Errorf("failed to upload file: %v", err)
	}

	logger.AppLogger.Info("Successfully uploaded file to S3", zap.String("bucket", bucket), zap.String("key", key))
	return nil
}

var cloudFrontClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetFileFromCloudFrontOrS3(svc *s3.S3, bucket, key string) (io.ReadCloser, error) {
	cloudfrontURL := fmt.Sprintf("%s/%s", appconst.AWSCloudFrontDomainName, key)
	resp, err := cloudFrontClient.Get(cloudfrontURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	file, err := GetS3File(svc, bucket, key)
	if err != nil {
		logger.AppLogger.Error("Failed to get file from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	return file, nil
}

func RemoveFileFromS3(svc *s3.S3, bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := svc.DeleteObject(input)
	if err != nil {
		logger.AppLogger.Error("Failed to delete object from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return fmt.Errorf("failed to delete object: %v", err)
	}

	logger.AppLogger.Info("Successfully deleted object from S3", zap.String("bucket", bucket), zap.String("key", key))
	return nil
}
