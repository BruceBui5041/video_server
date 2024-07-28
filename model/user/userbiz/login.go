package userbiz

import (
	"context"
	"video_server/common"
	"video_server/component/hasher"
	"video_server/component/tokenprovider"
	models "video_server/model"
	"video_server/model/user/usermodel"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

// type TokenConfig interface {
// 	GetAtExp() int
// 	GetRtExp() int
// }

type loginBusiness struct {
	loginStorage  LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        hasher.Hasher
	expiry        int
}

func NewLoginBusiness(
	storeUser LoginStorage,
	tokenProvicer tokenprovider.Provider,
	hasher hasher.Hasher,
	expiry int,
) *loginBusiness {
	return &loginBusiness{
		loginStorage:  storeUser,
		tokenProvider: tokenProvicer,
		hasher:        hasher,
		expiry:        expiry,
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

	// refreshToken, err := business.tokenProvider.Generate(payload, business.tokenConfig.GetRtExp())
	// if err != nil {
	// 	return nil, common.ErrInternal(err)
	// }

	// account := usermodel.NewAccount(accessToken, refreshToken)

	return accessToken, nil
}
