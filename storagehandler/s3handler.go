package storagehandler

import (
	"fmt"
	"io"
	"log"
	"os"
	"video_server/appconst"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var AWSSession *session.Session

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
		log.Fatal(err)
	}

	AWSSession = sess
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
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	fmt.Printf("Successfully retrieved %s from bucket %s\n", key, bucket)
	return result.Body, nil
}
