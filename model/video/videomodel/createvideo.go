package videomodel

import (
	"video_server/common"
)

type CreateVideo struct {
	common.SQLModel `json:",inline"`
	CourseID        uint   `json:"course_id" form:"course_id"`
	CourseSlug      string `json:"course_slug" form:"course_slug"`
	Title           string `json:"title" form:"title"`
	Slug            string `json:"slug" form:"slug"`
	Description     string `json:"description" form:"description"`
	VideoURL        string `json:"video_url" form:"video_url"`
	ThumbnailURL    string `json:"thumbnail_url" form:"thumbnail_url"`
	Duration        int    `json:"duration" form:"duration"`
	Order           int    `json:"order" form:"order"`
}

func (cv *CreateVideo) Mask(isAdmin bool) {
	cv.GenUID(common.DbTypeVideo)
}
