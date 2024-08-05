package usermodel

import models "video_server/model"

type CacheUser struct {
	Status            string            `json:"status"`
	LastName          string            `json:"lastname" gorm:"column:lastname;"`
	FirstName         string            `json:"firstname" gorm:"column:firstname;"`
	Email             string            `json:"email"`
	ProfilePictureURL string            `json:"profile_picture_url"`
	Roles             []models.Role     `json:"roles"`
	Auths             []models.UserAuth `json:"auths"`
}
