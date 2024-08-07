package videobiz

import (
	"context"
	"errors"
	"fmt"
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
	CheckExistingIntroVideo(ctx context.Context, courseID uint32) (*models.Video, error)
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

	if input.IntroVideo {
		existingIntroVideo, err := v.repo.CheckExistingIntroVideo(ctx, input.CourseID)
		if err != nil && err != common.RecordNotFound {
			return nil, common.ErrCannotCreateEntity(models.VideoEntityName, err)
		}
		if existingIntroVideo != nil {
			return nil, fmt.Errorf("the %s already is the intro video of the course", existingIntroVideo.Title)
		}
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
