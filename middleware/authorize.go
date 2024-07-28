package middleware

import (
	"errors"
	"strings"
	"video_server/common"
	"video_server/component"
	"video_server/component/tokenprovider/jwt"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err, "wrong authen header", "ErrWrongAuthHeader")
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(errors.New("wrong authen header"))
	}

	return parts[1], nil
}

func RequiredAuth(appCtx component.AppContext) func(ctx *gin.Context) {
	jwtProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
	return func(ctx *gin.Context) {
		token, err := extractTokenFromHeaderString(ctx.GetHeader("Authorization"))

		if err != nil {
			panic(err)
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

		if user.Status == 0 {
			panic(common.ErrNoPermission(errors.New("Account unavailable")))
		}

		user.Mask(false)

		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
