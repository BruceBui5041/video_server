package videobiz

import (
	"context"
	"errors"
	"io"
	"video_server/appconst"
	"video_server/common"
	models "video_server/model"
	"video_server/model/video/videomodel"
	"video_server/storagehandler"
)

type UpdateVideoRepo interface {
	UpdateVideo(ctx context.Context, id uint32, input *videomodel.UpdateVideo) (*models.Video, error)
	GetVideo(ctx context.Context, id uint32) (*models.Video, error)
}

type updateVideoBiz struct {
	repo UpdateVideoRepo
}

func NewUpdateVideoBiz(repo UpdateVideoRepo) *updateVideoBiz {
	return &updateVideoBiz{repo: repo}
}

func (v *updateVideoBiz) UpdateVideo(
	ctx context.Context,
	id uint32,
	input *videomodel.UpdateVideo,
	videoReader io.Reader,
	thumbnailReader io.Reader,
	userId uint32,
) (*models.Video, error) {
	// Validate input
	if err := v.validateInput(input); err != nil {
		return nil, err
	}

	// Get existing video
	existingVideo, err := v.repo.GetVideo(ctx, id)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.VideoEntityName, err)
	}

	// Handle file uploads
	if err := v.handleFileUploads(input, videoReader, thumbnailReader, userId, existingVideo.Id); err != nil {
		return nil, err
	}

	// Update video
	video, err := v.repo.UpdateVideo(ctx, id, input)
	if err != nil {
		return nil, common.ErrCannotUpdateEntity(models.VideoEntityName, err)
	}

	return video, nil
}

func (v *updateVideoBiz) validateInput(input *videomodel.UpdateVideo) error {
	if input.Title != nil && *input.Title == "" {
		return errors.New("video title cannot be empty")
	}

	if input.Slug != nil && *input.Slug == "" {
		return errors.New("video slug cannot be empty")
	}

	if input.Title != nil && len(*input.Title) > 255 {
		return errors.New("video title must not exceed 255 characters")
	}

	if input.Slug != nil && len(*input.Slug) > 255 {
		return errors.New("video slug must not exceed 255 characters")
	}

	return nil
}

func (v *updateVideoBiz) handleFileUploads(
	input *videomodel.UpdateVideo,
	videoReader io.Reader,
	thumbnailReader io.Reader,
	userId uint32,
	videoId uint32,
) error {
	fakeVideo := common.SQLModel{Id: videoId}
	fakeVideo.GenUID(common.DbTypeVideo)

	fakeUsr := common.SQLModel{Id: userId}
	fakeUsr.GenUID(common.DbTypeUser)

	if videoReader != nil {
		videoStorageInfo := storagehandler.VideoInfo{
			UploadedBy: fakeUsr.FakeId.String(),
			VideoId:    fakeVideo.FakeId.String(),
		}

		videoKey := storagehandler.GenerateVideoS3Key(videoStorageInfo)
		err := storagehandler.UploadFileToS3(videoReader, appconst.AWSVideoS3BuckerName, videoKey)
		if err != nil {
			return errors.New("failed to upload video to S3")
		}

		input.VideoURL = &videoKey
	}

	if thumbnailReader != nil {
		thumbnailStorageInfo := storagehandler.VideoInfo{
			UploadedBy: fakeUsr.FakeId.String(),
			VideoId:    fakeVideo.FakeId.String(),
		}

		thumbnailKey := storagehandler.GenerateThumbnailS3Key(thumbnailStorageInfo)
		err := storagehandler.UploadFileToS3(thumbnailReader, appconst.AWSVideoS3BuckerName, thumbnailKey)
		if err != nil {
			return errors.New("failed to upload thumbnail to S3")
		}

		input.ThumbnailURL = &thumbnailKey
	}

	return nil
}
