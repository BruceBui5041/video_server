package coursetransport

import (
	"net/http"
	"strconv"
	"video_server/common"
	"video_server/component"
	"video_server/model/category/categorystore"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursemodel"
	"video_server/model/course/courserepo"
	"video_server/model/course/coursestore"

	"github.com/gin-gonic/gin"
)

func UpdateCourseHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		var input coursemodel.UpdateCourse
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		categoryStore := categorystore.NewSQLStore(db)
		courseStore := coursestore.NewSQLStore(db)
		repo := courserepo.NewUpdateCourseRepo(courseStore, categoryStore)
		courseBusiness := coursebiz.NewUpdateCourseBiz(repo)

		course, err := courseBusiness.UpdateCourse(ctx, uint32(id), &input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(course))
	}
}
