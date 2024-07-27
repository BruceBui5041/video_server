package models

import "time"

type User struct {
	UserID            uint      `gorm:"primaryKey;autoIncrement"`
	Username          string    `gorm:"uniqueIndex;not null;size:50"`
	Email             string    `gorm:"uniqueIndex;not null;size:100"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	IsActive          bool      `gorm:"default:true"`
	ProfilePictureURL string    `gorm:"size:255"`
	Roles             []Role    `gorm:"many2many:user_roles;"`
	Auths             []UserAuth
	CreatedCourses    []Course `gorm:"foreignKey:CreatorID"`
	Enrollments       []Enrollment
	Progress          []Progress
}

func (User) TableName() string {
	return "user"
}
