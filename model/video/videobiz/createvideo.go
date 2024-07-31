package videobiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/video/videomodel"
)

type VideoRepo interface {
	CreateNewVideo(ctx context.Context, input *videomodel.CreateVideo) (*models.Video, error)
}

type createVideoBiz struct {
	repo VideoRepo
}

func NewCreateVideoBiz(repo VideoRepo) *createVideoBiz {
	return &createVideoBiz{repo: repo}
}

func (v *createVideoBiz) CreateNewVideo(ctx context.Context, input *videomodel.CreateVideo) (*models.Video, error) {
	if input.Title == "" {
		return nil, errors.New("video title is required")
	}

	if input.CourseID == 0 {
		return nil, errors.New("course ID is required")
	}

	if input.Slug == "" {
		return nil, errors.New("video slug is required")
	}

	if input.VideoURL == "" {
		return nil, errors.New("video URL is required")
	}

	if len(input.Title) > 255 {
		return nil, errors.New("video title must not exceed 255 characters")
	}

	if len(input.Slug) > 255 {
		return nil, errors.New("video slug must not exceed 255 characters")
	}

	video, err := v.repo.CreateNewVideo(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.VideoEntityName, err)
	}

	return video, nil
}
