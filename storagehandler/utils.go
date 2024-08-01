package storagehandler

import (
	"fmt"
	"path/filepath"
	"video_server/utils"
)

type VideoInfo struct {
	Useremail  string
	CourseSlug string
	VideoSlug  string
	Filename   string
}

func GenerateVideoS3Key(info VideoInfo) string {
	return fmt.Sprintf("course/%s/%s/%s/video_segment/%s",
		info.Useremail,
		info.CourseSlug,
		info.VideoSlug,
		utils.RenameFile(info.Filename, info.VideoSlug),
	)
}

func GenerateThumbnailS3Key(info VideoInfo) string {
	thumbnailFilename := generateThumbnailFilename(info.Filename)
	return fmt.Sprintf("course/%s/%s/%s/thumbnail/%s",
		info.Useremail,
		info.CourseSlug,
		info.VideoSlug,
		thumbnailFilename,
	)
}

func generateThumbnailFilename(videoFilename string) string {
	extension := filepath.Ext(videoFilename)
	baseFilename := videoFilename[:len(videoFilename)-len(extension)]
	return baseFilename + "_thumbnail" + extension
}
