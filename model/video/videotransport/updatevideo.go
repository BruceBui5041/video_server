package videotransport

import (
	"net/http"
	"strconv"
	"video_server/common"
	"video_server/component"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func UpdateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
			return
		}

		var input videomodel.UpdateVideo

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		videoFile, _ := c.FormFile("video")
		thumbnailFile, _ := c.FormFile("thumbnail")

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		useremail := requester.GetEmail()

		db := appCtx.GetMainDBConnection()
		store := videostore.NewSQLStore(db)
		repo := videorepo.NewUpdateVideoRepo(store)
		biz := videobiz.NewUpdateVideoBiz(repo)

		var videoReader, thumbnailReader interface{ Read([]byte) (int, error) }

		if videoFile != nil {
			videoReader, err = videoFile.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open video file"})
				return
			}
			defer videoReader.(interface{ Close() error }).Close()
		}

		if thumbnailFile != nil {
			thumbnailReader, err = thumbnailFile.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open thumbnail file"})
				return
			}
			defer thumbnailReader.(interface{ Close() error }).Close()
		}

		video, err := biz.UpdateVideo(c.Request.Context(), uint32(id), &input, videoReader, thumbnailReader, useremail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
	}
}
