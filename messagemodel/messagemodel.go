package messagemodel

// VideoInfo represents the information about a newly uploaded video
type VideoInfo struct {
	VideoID     string `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     uint   `json:"owner_id"`
	Timestamp   int64  `json:"timestamp"`
	S3Key       string `json:"s3key"`
	Slug        string `json:"slug"`
}
