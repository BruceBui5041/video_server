package userbiz

import (
	"context"
	"encoding/json"
	"fmt"
	"video_server/appconst"
	"video_server/common"
	"video_server/component/cache"
	"video_server/component/hasher"
	"video_server/component/tokenprovider"
	"video_server/logger"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

// type TokenConfig interface {
// 	GetAtExp() int
// 	GetRtExp() int
// }

type loginBusiness struct {
	loginStorage   LoginStorage
	tokenProvider  tokenprovider.Provider
	hasher         hasher.Hasher
	expiry         int
	dynamodbClient *cache.DynamoDBClient
}

func NewLoginBusiness(
	storeUser LoginStorage,
	tokenProvicer tokenprovider.Provider,
	hasher hasher.Hasher,
	expiry int,
	dynamodbClient *cache.DynamoDBClient,
) *loginBusiness {
	return &loginBusiness{
		loginStorage:   storeUser,
		tokenProvider:  tokenProvicer,
		hasher:         hasher,
		expiry:         expiry,
		dynamodbClient: dynamodbClient,
	}
}

func (business *loginBusiness) Login(ctx context.Context, data *usermodel.UserLogin) (*tokenprovider.Token, error) {
	user, err := business.loginStorage.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if err != nil {
		return nil, usermodel.ErrUsernameOrPasswordInvalid
	}

	pwdHashed := business.hasher.Hash(data.Password + user.Salt)
	if user.Password != pwdHashed {
		return nil, usermodel.ErrUsernameOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: int(user.Id),
		Roles:  user.Roles,
	}

	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	cacheKey := fmt.Sprintf("%s:%d", appconst.USER_PREFIX, user.Id)

	var cacheUser usermodel.CacheUser
	err = copier.Copy(&cacheUser, user)

	if err != nil {
		return nil, common.ErrInternal(err)
	}

	cacheUserJson, err := json.Marshal(cacheUser)

	if err == nil {
		err = business.dynamodbClient.Set(cacheKey, string(cacheUserJson))
		if err != nil {
			// Log the error, but continue with the request
			logger.AppLogger.Error("Failed to cache user", zap.Error(err))
		}
	}

	// refreshToken, err := business.tokenProvider.Generate(payload, business.tokenConfig.GetRtExp())
	// if err != nil {
	// 	return nil, common.ErrInternal(err)
	// }

	// account := usermodel.NewAccount(accessToken, refreshToken)

	return accessToken, nil
}
