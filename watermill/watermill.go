package watermill

import (
	"video_server/component"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MessageHandler func(appCtx component.AppContext, msg *message.Message)

func NewPubsubPublisher() *gochannel.GoChannel {
	return gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
}
