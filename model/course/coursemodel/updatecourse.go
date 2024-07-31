package coursemodel

import (
	"video_server/common"
)

type UpdateCourse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	CategoryID      string `json:"category_id"`
	Slug            string `json:"slug"`
}

func (uc *UpdateCourse) Mask(isAdmin bool) {
	uc.GenUID(common.DbTypeCourse)
}
