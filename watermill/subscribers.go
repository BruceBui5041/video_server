package watermill

import (
	"context"
	"log"
	"video_server/appconst"

	"github.com/ThreeDotsLabs/watermill/message"
)

func StartSubscribers() {
	// Define topics and their handlers
	topicHandlers := map[string]MessageHandler{
		appconst.TopicNewVideoUploaded: HandleNewVideoUpload,
		appconst.TopicVideoProcessed:   HandleVideoProcessed,
		// Add more topics and handlers as needed
	}

	// Subscribe to all topics
	messageChannels := make(map[string]<-chan *message.Message)
	for topic := range topicHandlers {
		messages, err := Publisher.Subscribe(context.Background(), topic)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
		}
		messageChannels[topic] = messages
	}

	// Process messages from all topics
	processMessages(messageChannels, topicHandlers)
}

func processMessages(messageChannels map[string]<-chan *message.Message, topicHandlers map[string]MessageHandler) {
	for {
		select {
		case msg := <-messageChannels[appconst.TopicNewVideoUploaded]:
			topicHandlers[appconst.TopicNewVideoUploaded](msg)
		case msg := <-messageChannels[appconst.TopicVideoProcessed]:
			topicHandlers[appconst.TopicVideoProcessed](msg)
			// Add more cases for additional topics
		}
	}
}
