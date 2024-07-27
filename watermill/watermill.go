package watermill

import (
	"video_server/common"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MessageHandler func(appCtx common.AppContext, msg *message.Message)

func NewPubsubPublisher() *gochannel.GoChannel {
	return gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
}
