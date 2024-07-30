package models

import "video_server/common"

const CourseEntityName = "Course"

type Course struct {
	common.SQLModel `json:",inline"`
	Title           string   `gorm:"not null;size:255"`
	Description     string   `gorm:"type:text"`
	CreatorID       uint32   `gorm:"index"`
	CategoryID      uint32   `gorm:"index"`
	Creator         User     `gorm:"constraint:OnDelete:SET NULL;"`
	Category        Category `gorm:"constraint:OnDelete:SET NULL;"`
	Videos          []Video
	Enrollments     []Enrollment
	Slug            string `gorm:"not null;size:255"`
}

func (Course) TableName() string {
	return "course"
}

func (u *Course) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeCourse)
}
