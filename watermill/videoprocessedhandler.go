package watermill

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

func HandleVideoProcessed(msg *message.Message) {
	videoID := string(msg.Payload)
	fmt.Printf("Video processed: %s\n", videoID)
	// Add your processing logic here
	// For example, you could update the video status in a database
	msg.Ack()
}
