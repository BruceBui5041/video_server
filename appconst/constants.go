package appconst

const (
	RawVideoS3Key = "raw-videos"
)

const (
	TopicNewVideoUploaded = "new_video_uploaded"
	TopicVideoProcessing  = "video_processing"
	TopicVideoProcessed   = "video_processed"
)

const (
	MaxConcurrentS3Push     = 50
	AWSVideoS3BuckerName    = "hls-video-segment"
	AWSRegion               = "ap-southeast-1"
	AWSCloudFrontDomainName = "d17cfikyg12m49.cloudfront.net"
)
