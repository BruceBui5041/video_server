package videomodel

import (
	"video_server/common"
)

type CreateVideo struct {
	common.SQLModel `json:",inline"`
	CourseID        uint   `json:"course_id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	VideoURL        string `json:"video_url"`
	Duration        int    `json:"duration"`
	Order           int    `json:"order"`
}

func (cv *CreateVideo) Mask(isAdmin bool) {
	cv.GenUID(common.DbTypeVideo)
}
