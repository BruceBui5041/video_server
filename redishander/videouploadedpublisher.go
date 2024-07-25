package redishander

import (
	"context"
	"encoding/json"
	"log"
	"video_server/appconst"

	"github.com/go-redis/redis/v8"
)

// VideoInfo represents the information about a newly uploaded video
type VideoInfo struct {
	VideoID     string `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UploadedBy  string `json:"uploaded_by"`
	Timestamp   int64  `json:"timestamp"`
}

// PublishNewVideoUploaded publishes a message to the Redis channel when a new video is uploaded
func PublishNewVideoUploaded(redisClient *redis.Client, videoInfo VideoInfo) error {
	ctx := context.Background()

	// Convert the VideoInfo struct to JSON
	payload, err := json.Marshal(videoInfo)
	if err != nil {
		log.Printf("Error marshaling video info: %v", err)
		return err
	}

	// Publish the message to the Redis channel
	err = redisClient.Publish(ctx, appconst.TopicNewVideoUploaded, payload).Err()
	if err != nil {
		log.Printf("Error publishing message to Redis: %v", err)
		return err
	}

	log.Printf("Published new video upload notification for video ID: %s", videoInfo.VideoID)
	return nil
}
