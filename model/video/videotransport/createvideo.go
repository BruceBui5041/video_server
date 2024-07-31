package videotransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func CreateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input videomodel.CreateVideo

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)
		videoStore := videostore.NewSQLStore(db)
		repo := videorepo.NewCreateVideoRepo(videoStore, courseStore)
		videoBusiness := videobiz.NewCreateVideoBiz(repo)

		video, err := videoBusiness.CreateNewVideo(ctx, &input)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
	}
}
