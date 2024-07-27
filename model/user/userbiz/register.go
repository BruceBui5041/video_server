package userbiz

import (
	"context"
	"errors"
	"regexp"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"golang.org/x/crypto/bcrypt"
)

type RegisterStorage interface {
	CreateNewUser(ctx context.Context, input *usermodel.CreateUser) (*models.User, error)
}

type registerBiz struct {
	registerStorage RegisterStorage
}

func NewRegisterBusiness(registerStorage RegisterStorage) *registerBiz {
	return &registerBiz{registerStorage: registerStorage}
}

// RegisterUser handles the registration of a new user
func (registerBiz *registerBiz) RegisterUser(ctx context.Context, newUser *usermodel.CreateUser) (*models.User, error) {
	// Validate required fields
	if newUser.Username == "" {
		return nil, errors.New("username is required")
	}
	if newUser.Email == "" {
		return nil, errors.New("email is required")
	}
	if newUser.AuthType == "" {
		return nil, errors.New("auth type is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(newUser.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate auth type and corresponding fields
	switch newUser.AuthType {
	case "password":
		if newUser.Password == "" {
			return nil, errors.New("password is required for password auth type")
		}
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		newUser.Password = string(hashedPassword)
	case "oauth":
		if newUser.AuthProviderID == "" || newUser.AuthProviderToken == "" {
			return nil, errors.New("auth provider ID and token are required for oauth auth type")
		}
	default:
		return nil, errors.New("invalid auth type")
	}

	createdUser, err := registerBiz.registerStorage.CreateNewUser(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
