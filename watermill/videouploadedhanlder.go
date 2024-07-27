package watermill

import (
	"encoding/json"
	"video_server/common"
	"video_server/grpcserver"
	"video_server/logger"
	"video_server/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func HandleNewVideoUpload(appCtx common.AppContext, msg *message.Message) {
	var videoInfo *messagemodel.VideoInfo
	err := json.Unmarshal(msg.Payload, &videoInfo)
	if err != nil {
		logger.AppLogger.Error("Cannot unmarshal message payload", zap.Any("payload", msg.Payload), zap.Error(err))
		return
	}

	go grpcserver.ProcessNewVideoRequest(appCtx, videoInfo)
	msg.Ack()
}
