package redishandler

import (
	"context"
	"encoding/json"
	"log"
	"video_server/appconst"
	"video_server/messagemodel"
)

// PublishNewVideoUploaded publishes a message to the Redis channel when a new video is uploaded
func PublishNewVideoUploaded(videoInfo messagemodel.VideoInfo) error {
	ctx := context.Background()

	// Convert the VideoInfo struct to JSON
	payload, err := json.Marshal(videoInfo)
	if err != nil {
		log.Printf("Error marshaling video info: %v", err)
		return err
	}

	// Publish the message to the Redis channel
	err = RedisClient.Publish(ctx, appconst.TopicNewVideoUploaded, payload).Err()
	if err != nil {
		log.Printf("Error publishing message to Redis: %v", err)
		return err
	}

	log.Printf("Published new video upload notification for video ID: %s", videoInfo.VideoID)
	return nil
}
