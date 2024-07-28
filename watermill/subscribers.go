package watermill

import (
	"context"
	"log"
	"video_server/appconst"
	"video_server/component"

	"github.com/ThreeDotsLabs/watermill/message"
)

func StartSubscribers(appCtx component.AppContext) {
	// Define topics and their handlers
	topicHandlers := map[string]MessageHandler{
		appconst.TopicNewVideoUploaded: HandleNewVideoUpload,
		appconst.TopicVideoProcessed:   HandleVideoProcessed,
		// Add more topics and handlers as needed
	}

	// Subscribe to all topics
	messageChannels := make(map[string]<-chan *message.Message)
	for topic := range topicHandlers {
		messages, err := appCtx.GetLocalPublisher().Subscribe(context.Background(), topic)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
		}
		messageChannels[topic] = messages
	}

	// Process messages from all topics
	processMessages(appCtx, messageChannels, topicHandlers)
}

func processMessages(
	appCtx component.AppContext,
	messageChannels map[string]<-chan *message.Message,
	topicHandlers map[string]MessageHandler,
) {
	for {
		select {
		case msg := <-messageChannels[appconst.TopicNewVideoUploaded]:
			topicHandlers[appconst.TopicNewVideoUploaded](appCtx, msg)
		case msg := <-messageChannels[appconst.TopicVideoProcessed]:
			topicHandlers[appconst.TopicVideoProcessed](appCtx, msg)
			// Add more cases for additional topics
		}
	}
}
