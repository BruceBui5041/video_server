package models

import (
	"video_server/common"
)

type Video struct {
	common.SQLModel `json:",inline"`
	CourseID        uint   `gorm:"index"`
	Title           string `gorm:"not null;size:255"`
	Slug            string `gorm:"not null;size:255"`
	Description     string
	VideoURL        string `gorm:"not null;size:255"`
	Duration        int
	Order           int
	Course          Course `gorm:"constraint:OnDelete:CASCADE;"`
	Tags            []Tag  `gorm:"many2many:video_tags;"`
	Progress        []Progress
}

func (Video) TableName() string {
	return "video"
}

func (u *Video) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeVideo)
}
