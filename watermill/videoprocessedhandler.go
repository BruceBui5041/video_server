package watermill

import (
	"fmt"
	"video_server/component"

	"github.com/ThreeDotsLabs/watermill/message"
)

func HandleVideoProcessed(appCtx component.AppContext, msg *message.Message) {
	videoID := string(msg.Payload)
	fmt.Printf("Video processed: %s\n", videoID)
	// Add your processing logic here
	// For example, you could update the video status in a database
	msg.Ack()
}
