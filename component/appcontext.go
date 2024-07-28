package component

import (
	pb "video_server/proto/video_service/video_service"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetLocalPublisher() *gochannel.GoChannel
	GetVideoProcessingClient() pb.VideoProcessingServiceClient
	SecretKey() string
}

type appCtx struct {
	db                 *gorm.DB
	localPublisher     *gochannel.GoChannel
	videoServiceClient pb.VideoProcessingServiceClient
	jwtSecretKey       string
}

func NewAppContext(
	db *gorm.DB,
	localPublisher *gochannel.GoChannel,
	videoServiceClient pb.VideoProcessingServiceClient,
	jwtSecretKey string,
) *appCtx {
	return &appCtx{
		db:                 db,
		localPublisher:     localPublisher,
		videoServiceClient: videoServiceClient,
		jwtSecretKey:       jwtSecretKey,
	}
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

func (ctx *appCtx) SecretKey() string {
	return ctx.jwtSecretKey
}
