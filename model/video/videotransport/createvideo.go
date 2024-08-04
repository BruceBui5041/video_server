package videotransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"
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

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)
		videoStore := videostore.NewSQLStore(db)
		repo := videorepo.NewCreateVideoRepo(videoStore, courseStore)
		biz := videobiz.NewCreateVideoBiz(repo)

		video, err := biz.CreateNewVideo(c.Request.Context(), &input, videoFile, thumbnailFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requester.Mask(false)
		videoUploadedInfo := &messagemodel.VideoInfo{
			RawVidS3Key: video.VideoURL,
			UploadedBy:  requester.GetFakeId(),
			CourseId:    video.Course.FakeId.String(),
			VideoId:     video.FakeId.String(),
		}

		err = watermill.PublishVideoUploadedEvent(appCtx, videoUploadedInfo)
		if err != nil {
			logger.AppLogger.Error("publish video uploaded event", zap.Error(err), zap.String("filename", input.Slug))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish uploaded video event"})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
	}
}
