package models

import (
	"time"
)

type Role struct {
	RoleID      uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex;not null;size:50"`
	Description string
	Users       []User `gorm:"many2many:user_roles;"`
}

type Permission struct {
	PermissionID uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"uniqueIndex;not null;size:50"`
	Description  string
}

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

type UserAuth struct {
	AuthID            uint   `gorm:"primaryKey;autoIncrement"`
	UserID            uint   `gorm:"index"`
	AuthType          string `gorm:"not null;size:20"`
	AuthProviderID    string `gorm:"size:255"`
	AuthProviderToken string
	User              User `gorm:"constraint:OnDelete:CASCADE;"`
}

type Category struct {
	CategoryID  uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;size:100"`
	Description string
	Courses     []Course
}

type Course struct {
	CourseID    uint   `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"not null;size:255"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatorID   uint      `gorm:"index"`
	CategoryID  uint      `gorm:"index"`
	Creator     User      `gorm:"constraint:OnDelete:SET NULL;"`
	Category    Category  `gorm:"constraint:OnDelete:SET NULL;"`
	Videos      []Video
	Enrollments []Enrollment
}

type Video struct {
	VideoID     uint   `gorm:"primaryKey;autoIncrement"`
	CourseID    uint   `gorm:"index"`
	Title       string `gorm:"not null;size:255"`
	Description string
	VideoURL    string `gorm:"not null;size:255"`
	Duration    int
	Order       int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Course      Course    `gorm:"constraint:OnDelete:CASCADE;"`
	Tags        []Tag     `gorm:"many2many:video_tags;"`
	Progress    []Progress
}

type Tag struct {
	TagID  uint    `gorm:"primaryKey;autoIncrement"`
	Name   string  `gorm:"uniqueIndex;not null;size:50"`
	Videos []Video `gorm:"many2many:video_tags;"`
}

type Enrollment struct {
	EnrollmentID uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"index"`
	CourseID     uint      `gorm:"index"`
	EnrolledAt   time.Time `gorm:"autoCreateTime"`
	User         User      `gorm:"constraint:OnDelete:CASCADE;"`
	Course       Course    `gorm:"constraint:OnDelete:CASCADE;"`
}

type Progress struct {
	ProgressID     uint      `gorm:"primaryKey;autoIncrement"`
	UserID         uint      `gorm:"index"`
	VideoID        uint      `gorm:"index"`
	WatchedSeconds int       `gorm:"default:0"`
	Completed      bool      `gorm:"default:false"`
	LastWatched    time.Time `gorm:"autoUpdateTime"`
	User           User      `gorm:"constraint:OnDelete:CASCADE;"`
	Video          Video     `gorm:"constraint:OnDelete:CASCADE;"`
}
