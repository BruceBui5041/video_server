package videorepo

import (
	"context"
	models "video_server/model"
	"video_server/model/video/videomodel"
)

type UpdateVideoStore interface {
	UpdateVideo(
		ctx context.Context,
		id uint32,
		updateData *videomodel.UpdateVideo,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Video, error)
}

type updateVideoRepo struct {
	store UpdateVideoStore
}

func NewUpdateVideoRepo(store UpdateVideoStore) *updateVideoRepo {
	return &updateVideoRepo{store: store}
}

func (repo *updateVideoRepo) UpdateVideo(ctx context.Context, id uint32, input *videomodel.UpdateVideo) (*models.Video, error) {
	if err := repo.store.UpdateVideo(ctx, id, input); err != nil {
		return nil, err
	}

	return repo.GetVideo(ctx, id)
}

func (repo *updateVideoRepo) GetVideo(ctx context.Context, id uint32) (*models.Video, error) {
	return repo.store.FindOne(ctx, map[string]interface{}{"id": id})
}
