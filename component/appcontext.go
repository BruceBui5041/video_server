package component

import (
	"video_server/component/cache"
	pb "video_server/proto/video_service/video_service"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetLocalPublisher() *gochannel.GoChannel
	GetVideoProcessingClient() pb.VideoProcessingServiceClient
	SecretKey() string
	GetDynamoDBClient() *cache.DynamoDBClient
}

type appCtx struct {
	db                 *gorm.DB
	localPublisher     *gochannel.GoChannel
	videoServiceClient pb.VideoProcessingServiceClient
	jwtSecretKey       string
	dynamoDBClient     *cache.DynamoDBClient
}

func NewAppContext(
	db *gorm.DB,
	localPublisher *gochannel.GoChannel,
	videoServiceClient pb.VideoProcessingServiceClient,
	jwtSecretKey string,
	dynamoDBClient *cache.DynamoDBClient,
) *appCtx {
	return &appCtx{
		db:                 db,
		localPublisher:     localPublisher,
		videoServiceClient: videoServiceClient,
		jwtSecretKey:       jwtSecretKey,
		dynamoDBClient:     dynamoDBClient,
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

func (ctx *appCtx) GetDynamoDBClient() *cache.DynamoDBClient {
	return ctx.dynamoDBClient
}
