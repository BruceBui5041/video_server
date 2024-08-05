package usertransport

import (
	"net/http"
	"time"
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
		loginbiz := userbiz.NewLoginBusiness(
			userStore,
			tokenProvider,
			md5,
			60*60*24*30,
			appCtx.GetDynamoDBClient(),
		)

		account, err := loginbiz.Login(ctx, &loginUserData)

		if err != nil {
			panic(err)
		}

		// Set the access token as a cookie
		cookie := &http.Cookie{
			Name:     "access_token",
			Value:    account.Token,
			HttpOnly: true,
			Secure:   false, // Set to true if using HTTPS
			// SameSite: http.SameSiteStrictMode,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Domain:   "localhost",
			Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
		}

		http.SetCookie(ctx.Writer, cookie)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
