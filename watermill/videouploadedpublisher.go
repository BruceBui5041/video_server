package watermill

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"video_server/appconst"
	"video_server/component"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/storagehandler"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func PublishVideoUploadedEvent(
	appCtx component.AppContext,
	videoInfo *messagemodel.VideoInfo,
	file multipart.File,
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

	err = storagehandler.UploadFileToS3(file, appconst.AWSVideoS3BuckerName, videoInfo.S3Key)
	if err != nil {
		logger.AppLogger.Error(
			"Failed to upload video to S3",
			zap.Any("videoInfo", videoInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(videoInfo.VideoID, payload)
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
