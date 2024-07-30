package models

import "video_server/common"

const CourseEntityName = "Course"

type Course struct {
	common.SQLModel `json:",inline"`
	Title           string       `json:"title" gorm:"not null;size:255"`
	Description     string       `json:"description" gorm:"type:text"`
	CreatorID       uint32       `json:"creator_id" gorm:"index"`
	CategoryID      uint32       `json:"category_id" gorm:"index"`
	Creator         User         `json:"creator" gorm:"constraint:OnDelete:SET NULL;"`
	Category        Category     `json:"category" gorm:"constraint:OnDelete:SET NULL;"`
	Videos          []Video      `json:"videos"`
	Enrollments     []Enrollment `json:"enrollments"`
	Slug            string       `json:"slug" gorm:"not null;size:255"`
}

func (Course) TableName() string {
	return "course"
}

func (u *Course) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeCourse)
}
