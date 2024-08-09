package middleware

import (
	"encoding/json"
	"errors"
	"video_server/common"
	"video_server/component"
	"video_server/component/tokenprovider/jwt"
	models "video_server/model"

	"github.com/gin-gonic/gin"
)

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err, "wrong authen header", "ErrWrongAuthHeader")
}

func RequiredAuth(appCtx component.AppContext) func(ctx *gin.Context) {
	jwtProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")

		if err != nil {
			panic(ErrWrongAuthHeader(errors.New("access_token cookie not found")))
		}

		payload, err := jwtProvider.Validate(token)

		if err != nil {
			panic(err)
		}

		// Try to get user from cache
		dynamoDBClient := appCtx.GetDynamoDBClient()
		cachedUser, err := dynamoDBClient.GetUserCache(uint32(payload.UserId))
		if err != nil || cachedUser == "" {
			panic(common.ErrNoPermission(errors.New("token expired")))
		}

		var user *models.User
		// User found in cache, unmarshal it
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			// If there's an error unmarshalling, we'll fetch from the database
			user = nil
		}

		if user.Status != common.StatusActive {
			panic(common.ErrNoPermission(errors.New("account unavailable")))
		}

		uid, err := common.FromBase58(user.GetFakeId())
		if err != nil {
			panic(err)
		}

		user.Id = uid.GetLocalID()

		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
