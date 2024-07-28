package coursemodel

import (
	"video_server/common"
)

// CreateUser represents the data needed to create a new user
type CreateCourse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	CategoryID      uint   `json:"category_id"`
	CreatorID       uint   `json:"creator_id"`
	Slug            string `json:"slug"`
}

func (cc *CreateCourse) Mask(isAdmin bool) {
	cc.GenUID(common.DbTypeCourse)
}
