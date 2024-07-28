package usertransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/component/hasher"
	"video_server/model/user/userbiz"
	"video_server/model/user/usermodel"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func Register(appCtx component.AppContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		db := appCtx.GetMainDBConnection()
		var data usermodel.CreateUser

		if err := ctx.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := userstore.NewSQLStore(db)
		md5 := hasher.NewMD5Hash()
		business := userbiz.NewRegisterBusiness(store, md5)

		if err := business.RegisterUser(ctx, &data); err != nil {
			panic(err)
		}

		data.Mask(false)

		ctx.JSON(http.StatusCreated, common.SimpleSuccessResponse(data.FakeId.String()))

	}
}
