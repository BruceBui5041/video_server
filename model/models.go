package models

import (
	"time"
)

type Permission struct {
	PermissionID uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"uniqueIndex;not null;size:50"`
	Description  string
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
