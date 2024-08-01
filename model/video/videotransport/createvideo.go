package videotransport

import (
	"net/http"
	"video_server/appconst"
	"video_server/common"
	"video_server/component"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"
	"video_server/storagehandler"
	"video_server/watermill"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input videomodel.CreateVideo

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		videoFile, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No video file uploaded"})
			return
		}

		thumbnailFile, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No thumbnail file uploaded"})
			return
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		useremail := requester.GetEmail()

		videoStorageInfo := storagehandler.VideoInfo{
			Useremail:  useremail,
			CourseSlug: input.CourseSlug,
			VideoSlug:  input.Slug,
			Filename:   videoFile.Filename,
		}

		videoKey := storagehandler.GenerateVideoS3Key(videoStorageInfo)
		thumbnailKey := storagehandler.GenerateThumbnailS3Key(videoStorageInfo)

		// Open video file
		videoFileContent, err := videoFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open video file"})
			return
		}
		defer videoFileContent.Close()

		// Upload video to S3
		err = storagehandler.UploadFileToS3(videoFileContent, appconst.AWSVideoS3BuckerName, videoKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video to S3"})
			return
		}

		// Open thumbnail file
		thumbnailFileContent, err := thumbnailFile.Open()
		if err != nil {
			// Remove the uploaded video if thumbnail file opening fails
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, videoKey)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open thumbnail file"})
			return
		}
		defer thumbnailFileContent.Close()

		// Upload thumbnail to S3
		err = storagehandler.UploadFileToS3(thumbnailFileContent, appconst.AWSVideoS3BuckerName, thumbnailKey)
		if err != nil {
			// Remove the uploaded video if thumbnail upload fails
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, videoKey)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload thumbnail to S3"})
			return
		}

		// Update input with S3 URLs
		input.VideoURL = videoKey
		input.ThumbnailURL = thumbnailKey

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)
		videoStore := videostore.NewSQLStore(db)
		repo := videorepo.NewCreateVideoRepo(videoStore, courseStore)
		biz := videobiz.NewCreateVideoBiz(repo)

		video, err := biz.CreateNewVideo(c.Request.Context(), &input)
		if err != nil {
			// Remove both video and thumbnail from S3 if video creation fails
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, videoKey)
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, thumbnailKey)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		videoUploadedInfo := &messagemodel.VideoInfo{
			RawVidS3Key: videoKey,
			CourseSlug:  input.CourseSlug,
			VideoSlug:   input.Slug,
			UserEmail:   useremail,
		}

		err = watermill.PublishVideoUploadedEvent(appCtx, videoUploadedInfo)
		if err != nil {
			logger.AppLogger.Error("publish video uploaded event", zap.Error(err), zap.String("filename", input.Slug))
			// Remove both video and thumbnail from S3 if video creation fails
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, videoKey)
			_ = storagehandler.RemoveFileFromS3(appconst.AWSVideoS3BuckerName, thumbnailKey)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish uploaded video event"})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
	}
}
