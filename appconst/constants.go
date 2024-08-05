package appconst

// DynamoDB
const (
	VIDEO_URL_PREFIX = "video_url"
	USER_PREFIX      = "user"
)

// local TOPIC
const (
	TopicNewVideoUploaded = "new_video_uploaded"
	TopicVideoProcessing  = "video_processing"
	TopicVideoProcessed   = "video_processed"
)

const (
	MaxConcurrentS3Push     = 50
	AWSVideoS3BuckerName    = "hls-video-segment"
	AWSCloudFrontDomainName = "https://d17cfikyg12m49.cloudfront.net"
)
