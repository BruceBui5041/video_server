package watermill

import (
	"encoding/json"
	"video_server/grpcserver"
	"video_server/logger"
	"video_server/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func HandleNewVideoUpload(msg *message.Message) {
	var videoInfo *messagemodel.VideoInfo
	err := json.Unmarshal(msg.Payload, videoInfo)
	if err != nil {
		logger.AppLogger.Error("Cannot unmarshal message payload", zap.Any("payload", msg.Payload))
		return
	}

	go grpcserver.RequestNewVideoUploaded(videoInfo)
	msg.Ack()
}
