package common

import (
	pb "video_server/proto/video_service/video_service"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetLocalPublisher() *gochannel.GoChannel
	GetVideoProcessingClient() pb.VideoProcessingServiceClient
}

type appCtx struct {
	db                 *gorm.DB
	localPublisher     *gochannel.GoChannel
	videoServiceClient pb.VideoProcessingServiceClient
}

func NewAppContext(
	db *gorm.DB,
	localPublisher *gochannel.GoChannel,
	videoServiceClient pb.VideoProcessingServiceClient,
) *appCtx {
	return &appCtx{db: db, localPublisher: localPublisher, videoServiceClient: videoServiceClient}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}

func (ctx *appCtx) GetLocalPublisher() *gochannel.GoChannel {
	return ctx.localPublisher
}

func (ctx *appCtx) GetVideoProcessingClient() pb.VideoProcessingServiceClient {
	return ctx.videoServiceClient
}
