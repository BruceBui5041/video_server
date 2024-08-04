package videobiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type GetVideoRepo interface {
	GetVideo(ctx context.Context, id uint32, courseSlug string) (*models.Video, error)
}

type getVideoBiz struct {
	repo GetVideoRepo
}

func NewGetVideoBiz(repo GetVideoRepo) *getVideoBiz {
	return &getVideoBiz{repo: repo}
}

func (biz *getVideoBiz) GetVideoById(ctx context.Context, id uint32, courseSlug string) (*models.Video, error) {
	video, err := biz.repo.GetVideo(ctx, id, courseSlug)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.VideoEntityName, err)
	}
	return video, nil
}
