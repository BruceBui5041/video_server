package watermill

import (
	"encoding/json"
	"fmt"
	"video_server/appconst"
	"video_server/component"
	"video_server/logger"
	"video_server/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func PublishVideoUploadedEvent(
	appCtx component.AppContext,
	videoInfo *messagemodel.VideoInfo,
) error {
	// Marshal videoInfo into JSON
	payload, err := json.Marshal(videoInfo)
	if err != nil {
		logger.AppLogger.Error(
			"Error marshaling videoInfo to JSON",
			zap.Any("videoInfo", videoInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(videoInfo.VideoSlug, payload)
	err = appCtx.GetLocalPublisher().Publish(appconst.TopicNewVideoUploaded, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(
			fmt.Sprintf("Error publish %s", appconst.TopicNewVideoUploaded),
			zap.Any("msg payload", payload),
			zap.Error(err),
		)
		return err
	}

	return nil
}
