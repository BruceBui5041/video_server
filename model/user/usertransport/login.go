package usertransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/component/hasher"
	"video_server/component/tokenprovider/jwt"
	"video_server/model/user/userbiz"
	"video_server/model/user/usermodel"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func Login(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := ctx.ShouldBind(&loginUserData); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		md5 := hasher.NewMD5Hash()

		userStore := userstore.NewSQLStore(db)
		loginbiz := userbiz.NewLoginBusiness(userStore, tokenProvider, md5, 60*60*24*30)

		account, err := loginbiz.Login(ctx, &loginUserData)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
