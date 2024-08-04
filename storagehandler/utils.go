package storagehandler

import (
	"fmt"
	"path/filepath"
)

type VideoInfo struct {
	UploadedBy        string `json:"uploaded_by"`
	CourseId          string `json:"course_id"`
	VideoId           string `json:"video_id"`
	ThumbnailFilename string `json:"thumbnail_filename"`
}

func GenerateVideoS3Key(info VideoInfo) string {
	return fmt.Sprintf("course/%s/%s/%s/video_segment/%s",
		info.UploadedBy,
		info.CourseId,
		info.VideoId,
		info.VideoId,
	)
}

func GenerateThumbnailS3Key(info VideoInfo) string {
	thumbnailFilename := generateThumbnailFilename(info)
	return fmt.Sprintf("course/%s/%s/%s/thumbnail/%s",
		info.UploadedBy,
		info.CourseId,
		info.VideoId,
		thumbnailFilename,
	)
}

func generateThumbnailFilename(info VideoInfo) string {
	videoFilename := info.ThumbnailFilename
	extension := filepath.Ext(videoFilename)
	return "thumbnail" + extension
}
