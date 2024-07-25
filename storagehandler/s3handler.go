package storagehandler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"video_server/appconst"
	"video_server/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var AWSSession *session.Session

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		logger.AppLogger.Fatal("Error loading .env file", zap.Error(err))
	}

	// Get credentials from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Create a new credential provider
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")

	// Create a new session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(appconst.AWSRegion),
		Credentials: creds,
	})
	if err != nil {
		logger.AppLogger.Fatal("Failed to create AWS session", zap.Error(err))
	}

	AWSSession = sess
	logger.AppLogger.Info("AWS session created successfully")
}

func GetS3File(bucket, key string) (io.ReadCloser, error) {
	// Create an S3 service client
	svc := s3.New(AWSSession)

	// Create the GetObject request
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Fetch the object
	result, err := svc.GetObject(input)
	if err != nil {
		logger.AppLogger.Error("Failed to get object from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	logger.AppLogger.Info("Successfully retrieved object from S3", zap.String("bucket", bucket), zap.String("key", key))
	return result.Body, nil
}

func UploadFileToS3(file io.Reader, bucket, key string) error {
	// Create an S3 service client
	svc := s3.New(AWSSession)

	// Read the entire file into a buffer
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		logger.AppLogger.Error("Failed to read file into buffer", zap.Error(err))
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Create the PutObject request
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

func GetFileFromCloudFrontOrS3(bucket, key string) (io.ReadCloser, error) {
	cloudfrontURL := fmt.Sprintf("%s/%s", appconst.AWSCloudFrontDomainName, key)
	resp, err := cloudFrontClient.Get(cloudfrontURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	file, err := GetS3File(bucket, key)
	if err != nil {
		logger.AppLogger.Error("Failed to get file from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	return file, nil
}
