package videorepo

import (
	"context"
	models "video_server/model"
	"video_server/model/video/videomodel"
)

type CourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Course, error)
}

type CreateVideoStore interface {
	CreateNewVideo(
		ctx context.Context,
		newVideo *models.Video,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Video, error)
}

type createVideoRepo struct {
	videoStore  CreateVideoStore
	courseStore CourseStore
}

func NewCreateVideoRepo(videoStore CreateVideoStore, courseStore CourseStore) *createVideoRepo {
	return &createVideoRepo{videoStore: videoStore, courseStore: courseStore}
}

func (repo *createVideoRepo) CreateNewVideo(ctx context.Context, input *videomodel.CreateVideo) (*models.Video, error) {
	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"slug": input.CourseSlug})
	if err != nil {
		return nil, err
	}

	newVideo := &models.Video{
		CourseID:     uint(course.Id),
		Title:        input.Title,
		Slug:         input.Slug,
		Description:  input.Description,
		VideoURL:     input.VideoURL,
		Duration:     input.Duration,
		Order:        input.Order,
		ThumbnailURL: input.ThumbnailURL,
	}

	videoId, err := repo.videoStore.CreateNewVideo(ctx, newVideo)
	if err != nil {
		return nil, err
	}

	video, err := repo.videoStore.FindOne(ctx, map[string]interface{}{"id": videoId})
	if err != nil {
		return nil, err
	}

	video.Course = *course

	return video, nil
}
