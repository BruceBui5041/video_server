package user

import (
	"errors"
	"time"
	models "video_server/model"

	"gorm.io/gorm"
)

// CreateUser represents the data needed to create a new user
type CreateUser struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	AuthType          string `json:"auth_type"`
	AuthProviderID    string `json:"auth_provider_id,omitempty"`
	AuthProviderToken string `json:"auth_provider_token,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}

// CreateUser creates a new user with the specified authentication type
func CreateNewUser(db *gorm.DB, input CreateUser) (*models.User, error) {
	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Check if user already exists
	var existingUser models.User
	if err := tx.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("user with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, err
	}

	// Create new user
	newUser := models.User{
		Username:          input.Username,
		Email:             input.Email,
		ProfilePictureURL: input.ProfilePictureURL,
		CreatedAt:         time.Now(),
		IsActive:          true,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create user authentication entry
	auth := models.UserAuth{
		UserID:            newUser.UserID,
		AuthType:          input.AuthType,
		AuthProviderID:    input.AuthProviderID,
		AuthProviderToken: input.AuthProviderToken,
	}

	// If it's a local auth type, we need to hash the password
	if input.AuthType == "local" {
		if input.Password == "" {
			tx.Rollback()
			return nil, errors.New("password is required for local authentication")
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
		return nil, err
	}

	// Assign default role (assuming 'user' role exists with ID 1)
	if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", newUser.UserID, 1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}
