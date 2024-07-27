package userbiz

import (
	"context"
	"video_server/common"
	"video_server/model/user/usermodel"

	"golang.org/x/crypto/bcrypt"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}) (*usermodel.User, error)
}

type loginBusiness struct {
	storeUser     LoginStorage
	hasher        Hasher
	tokenProvider TokenProvider
}

func NewLoginBusiness(storeUser LoginStorage, hasher Hasher, tokenProvider TokenProvider) *loginBusiness {
	return &loginBusiness{
		storeUser:     storeUser,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (business *loginBusiness) Login(ctx context.Context, data *usermodel.UserLogin) (*common.Token, error) {
	user, err := business.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if err != nil {
		return nil, common.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	passHashed := user.Password

	if err := business.hasher.ComparePassword(passHashed, data.Password); err != nil {
		return nil, common.ErrEmailOrPasswordInvalid
	}

	payload := common.TokenPayload{
		UserId: user.Id,
		Role:   user.Role,
	}

	accessToken, err := business.tokenProvider.Generate(payload, common.AccessTokenExpiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	refreshToken, err := business.tokenProvider.Generate(payload, common.RefreshTokenExpiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	account := common.NewAccount(accessToken, refreshToken)

	return account, nil
}

type Hasher interface {
	ComparePassword(hashedPassword, password string) error
}

type TokenProvider interface {
	Generate(data common.TokenPayload, expiry int) (*common.Token, error)
}

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
