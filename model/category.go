package models

import "video_server/common"

const CategoryEntityName = "Category"

type Category struct {
	common.SQLModel `json:",inline"`
	Name            string   `json:"name" gorm:"not null;size:100"`
	Description     string   `json:"description"`
	Courses         []Course `json:"course"`
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCategory)
}
