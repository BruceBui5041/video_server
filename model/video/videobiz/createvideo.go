package videobiz

import (
	"context"
	"errors"
	"mime/multipart"
	"video_server/common"
	models "video_server/model"
	"video_server/model/video/videomodel"
)

type VideoRepo interface {
	CreateNewVideo(
		ctx context.Context,
		input *videomodel.CreateVideo,
		videoFile,
		thumbnailFile *multipart.FileHeader,
	) (*models.Video, error)
}

type createVideoBiz struct {
	repo VideoRepo
}

func NewCreateVideoBiz(repo VideoRepo) *createVideoBiz {
	return &createVideoBiz{repo: repo}
}

func (v *createVideoBiz) CreateNewVideo(
	ctx context.Context,
	input *videomodel.CreateVideo,
	videoFile,
	thumbnailFile *multipart.FileHeader,
) (*models.Video, error) {
	if err := v.validateInput(input); err != nil {
		return nil, err
	}

	video, err := v.repo.CreateNewVideo(
		ctx,
		input,
		videoFile,
		thumbnailFile,
	)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.VideoEntityName, err)
	}

	return video, nil
}

func (v *createVideoBiz) validateInput(input *videomodel.CreateVideo) error {
	if input.Title == "" {
		return errors.New("video title is required")
	}
	if input.CourseSlug == "" {
		return errors.New("course slug is required")
	}
	if input.Slug == "" {
		return errors.New("video slug is required")
	}
	if len(input.Title) > 255 {
		return errors.New("video title must not exceed 255 characters")
	}
	if len(input.Slug) > 255 {
		return errors.New("video slug must not exceed 255 characters")
	}
	return nil
}
