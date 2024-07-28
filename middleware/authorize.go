package middleware

import (
	"errors"
	"video_server/common"
	"video_server/component"
	"video_server/component/tokenprovider/jwt"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err, "wrong authen header", "ErrWrongAuthHeader")
}

func RequiredAuth(appCtx component.AppContext) func(ctx *gin.Context) {
	jwtProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
	return func(ctx *gin.Context) {
		// Get the access_token from the cookie
		token, err := ctx.Cookie("access_token")

		if err != nil {
			panic(ErrWrongAuthHeader(errors.New("access_token cookie not found")))
		}

		payload, err := jwtProvider.Validate(token)

		if err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		userStore := userstore.NewSQLStore(db)

		user, err := userStore.FindUser(ctx, map[string]interface{}{"id": payload.UserId})
		if err != nil {
			panic(err)
		}

		if user.Status != common.StatusActive {
			panic(common.ErrNoPermission(errors.New("account unavailable")))
		}

		user.Mask(false)

		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
