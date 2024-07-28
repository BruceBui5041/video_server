package userstore

import (
	"context"
	"errors"
	models "video_server/model"
	user "video_server/model/user/usermodel"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewUser(ctx context.Context, input *user.CreateUser) error {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Check if user already exists
	var existingUser models.User
	if err := tx.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return errors.New("user with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	// Create new user
	newUser := models.User{
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Email:             input.Email,
		ProfilePictureURL: input.ProfilePictureURL,
		Password:          input.Password,
		Salt:              input.Salt,
		IsActive:          true,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create user authentication entry
	auth := models.UserAuth{
		UserID:            newUser.Id,
		AuthType:          input.AuthType,
		AuthProviderID:    input.AuthProviderID,
		AuthProviderToken: input.AuthProviderToken,
	}

	// If it's a local auth type, we need to hash the password
	if input.AuthType == "password" {
		if input.Password == "" {
			tx.Rollback()
			return errors.New("password is required for local authentication")
		}
		// TODO: Implement password hashing
		// hashedPassword, err := hashPassword(input.Password)
		// if err != nil {
		//     tx.Rollback()
		//     return nil, err
		// }
		// auth.AuthProviderToken = hashedPassword
	}

	if err := tx.Create(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Assign default role (assuming 'user' role exists with ID 1)
	if err := tx.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", newUser.Id, 1).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
