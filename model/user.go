package models

import (
	"encoding/json"
	"video_server/common"
	"video_server/logger"

	"go.uber.org/zap"
)

const UserEntityName = "User"

type User struct {
	common.SQLModel   `json:",inline"`
	LastName          string       `json:"lastname" gorm:"column:lastname;"`
	FirstName         string       `json:"firstname" gorm:"column:firstname;"`
	Email             string       `json:"email" gorm:"column:email;uniqueIndex;not null;size:100"`
	ProfilePictureURL string       `json:"profile_picture_url" gorm:"column:profile_picture_url;size:255"`
	Roles             []Role       `json:"roles" gorm:"many2many:user_role;"`
	Auths             []UserAuth   `json:"auths" gorm:"foreignKey:UserID"`
	CreatedCourses    []Course     `json:"created_courses" gorm:"foreignKey:CreatorID"`
	Enrollments       []Enrollment `json:"enrollments" `
	Progress          []Progress   `json:"progress"`
	Salt              string       `json:"-" gorm:"column:salt;"`
	Password          string       `json:"-" gorm:"column:password;"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

func (u *User) GetUserId() uint32 {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	data, err := json.Marshal(u.Roles)
	if err != nil {
		logger.AppLogger.Error("cannot marshal users roles", zap.Error(err))
		return ""
	}
	return string(data)
}

func (u *User) GetFakeId() string {
	return u.FakeId.String()
}
