package models

import (
	"video_server/common"
)

type User struct {
	common.SQLModel   `json:",inline"`
	LastName          string     `json:"lastname" gorm:"column:lastname;"`
	FirstName         string     `json:"firstname" gorm:"column:firstname;"`
	Email             string     `gorm:"column:email;uniqueIndex;not null;size:100"`
	ProfilePictureURL string     `gorm:"column:profile_picture_url;size:255"`
	Roles             []Role     `gorm:"many2many:user_role;"`
	Auths             []UserAuth `gorm:"foreignKey:UserID"`
	CreatedCourses    []Course   `gorm:"foreignKey:CreatorID"`
	Enrollments       []Enrollment
	Progress          []Progress
	Salt              string `json:"-" gorm:"column:salt;"`
	Password          string `json:"-" gorm:"column:password;"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}
