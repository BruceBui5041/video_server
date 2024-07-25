package watermill

import (
	"encoding/json"
	"fmt"
	"video_server/appconst"
	"video_server/logger"
	"video_server/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func PublishVideoUploaderEvent(videoInfo *messagemodel.VideoInfo) error {
	var msg message.Message
	err := json.Unmarshal([]byte(msg.Payload), &videoInfo)
	if err != nil {
		logger.AppLogger.Error(
			"Error parsing message payload",
			zap.Any("videoInfo", videoInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(videoInfo.VideoID, []byte(msg.Payload))
	err = Publisher.Publish(appconst.TopicNewVideoUploaded, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(
			fmt.Sprintf("Error publish %s", appconst.TopicNewVideoUploaded),
			zap.Any("msg payload", msg.Payload),
			zap.Error(err),
		)
		return err
	}

	return nil
}
