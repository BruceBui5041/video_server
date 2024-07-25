package redishandler

import (
	"context"
	"encoding/json"
	"log"
	"video_server/appconst"

	"github.com/go-redis/redis/v8"
)

func StartRedisSubscribers(redisClient *redis.Client) {
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, appconst.TopicNewVideoUploaded)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Received message from Redis channel %s: %s", msg.Channel, msg.Payload)

		// Parse the message payload
		var videoInfo struct {
			VideoID string `json:"video_id"`
			// Add other fields as needed
		}
		err := json.Unmarshal([]byte(msg.Payload), &videoInfo)
		if err != nil {
			log.Printf("Error parsing message payload: %v", err)
			continue
		}

		// // Create a Watermill message
		// watermillMsg := message.NewMessage(videoInfo.VideoID, []byte(msg.Payload))

		// // Process the message using the existing handler
		// watermill.HandleNewVideoUpload(watermillMsg)
	}
}
