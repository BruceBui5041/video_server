package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursemodel"
	"video_server/model/course/coursestore"

	"github.com/gin-gonic/gin"
)

func CreateCourseHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input coursemodel.CreateCourse

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		coursestore := coursestore.NewSQLStore(db)
		coursebusiness := coursebiz.NewCreateCourseBiz(coursestore)

		input.CreatorID = uint(requester.GetUserId())
		err := coursebusiness.CreateNewCourse(ctx, &input)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("ok"))
	}
}
