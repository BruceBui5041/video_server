package userbiz

import (
	"context"
	"errors"
	"regexp"
	"video_server/common"
	"video_server/component/hasher"
	"video_server/model/user/usermodel"
)

type RegisterStorage interface {
	CreateNewUser(ctx context.Context, input *usermodel.CreateUser) error
}

type registerBiz struct {
	registerStorage RegisterStorage
	hasher          hasher.Hasher
}

func NewRegisterBusiness(registerStorage RegisterStorage, hasher hasher.Hasher) *registerBiz {
	return &registerBiz{registerStorage: registerStorage, hasher: hasher}
}

// RegisterUser handles the registration of a new user
func (registerBiz *registerBiz) RegisterUser(ctx context.Context, inputData *usermodel.CreateUser) error {
	// Validate required fields
	if inputData.FirstName == "" {
		return errors.New("first name is required")
	}
	if inputData.LastName == "" {
		return errors.New("last name is required")
	}

	if inputData.Email == "" {
		return errors.New("email is required")
	}

	if inputData.Password != "" {
		inputData.AuthType = "password"
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(inputData.Email) {
		return errors.New("invalid email format")
	}

	// Validate auth type and corresponding fields
	switch inputData.AuthType {
	case "password":
		if inputData.Password == "" {
			return errors.New("password is required for password auth type")
		}
		// Hash password
		salt := common.GenSalt(50)
		hashedPassword := registerBiz.hasher.Hash(inputData.Password + salt)
		inputData.Salt = salt
		inputData.Password = string(hashedPassword)
	case "oauth":
		if inputData.AuthProviderID == "" || inputData.AuthProviderToken == "" {
			return errors.New("auth provider ID and token are required for oauth auth type")
		}
	default:
		return errors.New("invalid auth type")
	}

	err := registerBiz.registerStorage.CreateNewUser(ctx, inputData)

	if err != nil {
		return err
	}

	return nil
}
