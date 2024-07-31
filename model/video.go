package models

import (
	"video_server/common"
)

type Video struct {
	common.SQLModel `json:",inline"`
	CourseID        uint       `json:"course_id" gorm:"index"`
	Title           string     `json:"title" gorm:"not null;size:255"`
	Slug            string     `json:"slug" gorm:"not null;size:255"`
	Description     string     `json:"description"`
	VideoURL        string     `json:"video_url" gorm:"not null;size:255"`
	Duration        int        `json:"duration"`
	Order           int        `json:"order"`
	Course          Course     `json:"course" gorm:"constraint:OnDelete:CASCADE;"`
	Tags            []Tag      `json:"tags" gorm:"many2many:video_tags;"`
	Progress        []Progress `json:"progress"`
}

func (Video) TableName() string {
	return "video"
}

func (u *Video) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeVideo)
}
