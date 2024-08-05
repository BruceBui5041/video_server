package videorepo

import (
	"context"
	"errors"
	"mime/multipart"
	"video_server/appconst"
	"video_server/common"
	models "video_server/model"
	"video_server/model/video/videomodel"
	"video_server/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
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
	UpdateVideo(
		ctx context.Context,
		id uint32,
		updateData *videomodel.UpdateVideo,
	) error
}

type createVideoRepo struct {
	videoStore  CreateVideoStore
	courseStore CourseStore
	svc         *s3.S3
}

func NewCreateVideoRepo(videoStore CreateVideoStore, courseStore CourseStore, svc *s3.S3) *createVideoRepo {
	return &createVideoRepo{
		videoStore:  videoStore,
		courseStore: courseStore,
		svc:         svc,
	}
}

func (repo *createVideoRepo) CreateNewVideo(
	ctx context.Context,
	input *videomodel.CreateVideo,
	videoFile,
	thumbnailFile *multipart.FileHeader,
) (*models.Video, error) {
	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"slug": input.CourseSlug})
	if err != nil {
		return nil, err
	}

	newVideo := &models.Video{
		CourseID:     course.Id,
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

	video.Mask(false)
	course.Mask(false)

	sqlObj := common.SQLModel{Id: course.CreatorID}
	sqlObj.GenUID(common.DbTypeUser)

	videoStorageInfo := storagehandler.VideoInfo{
		UploadedBy:        sqlObj.FakeId.String(),
		CourseId:          course.FakeId.String(),
		VideoId:           video.FakeId.String(),
		ThumbnailFilename: thumbnailFile.Filename,
	}

	videoKey := storagehandler.GenerateVideoS3Key(videoStorageInfo)
	thumbnailKey := storagehandler.GenerateThumbnailS3Key(videoStorageInfo)

	if err := repo.uploadFiles(videoFile, thumbnailFile, videoKey, thumbnailKey); err != nil {
		// If video creation fails, remove the uploaded files
		repo.removeFiles(videoKey, thumbnailKey)
		return nil, err
	}

	video.VideoURL = videoKey
	video.ThumbnailURL = thumbnailKey

	err = repo.videoStore.UpdateVideo(
		ctx,
		videoId,
		&videomodel.UpdateVideo{VideoURL: &videoKey, ThumbnailURL: &thumbnailKey},
	)

	if err != nil {
		// If video creation fails, remove the uploaded files
		repo.removeFiles(videoKey, thumbnailKey)
		return nil, err
	}

	video.Course = *course

	return video, nil
}

func (repo *createVideoRepo) uploadFiles(videoFile, thumbnailFile *multipart.FileHeader, videoKey, thumbnailKey string) error {
	videoFileContent, err := videoFile.Open()
	if err != nil {
		return errors.New("failed to open video file")
	}
	defer videoFileContent.Close()

	err = storagehandler.UploadFileToS3(
		repo.svc,
		videoFileContent,
		appconst.AWSVideoS3BuckerName,
		videoKey,
	)
	if err != nil {
		return errors.New("failed to upload video to S3")
	}

	thumbnailFileContent, err := thumbnailFile.Open()
	if err != nil {
		go repo.removeFiles(videoKey, "")
		return errors.New("failed to open thumbnail file")
	}
	defer thumbnailFileContent.Close()

	err = storagehandler.UploadFileToS3(
		repo.svc,
		thumbnailFileContent,
		appconst.AWSVideoS3BuckerName,
		thumbnailKey,
	)
	if err != nil {
		go repo.removeFiles(videoKey, "")
		return errors.New("failed to upload thumbnail to S3")
	}

	return nil
}

func (repo *createVideoRepo) removeFiles(videoKey, thumbnailKey string) {
	if videoKey != "" {
		_ = storagehandler.RemoveFileFromS3(repo.svc, appconst.AWSVideoS3BuckerName, videoKey)
	}
	if thumbnailKey != "" {
		_ = storagehandler.RemoveFileFromS3(repo.svc, appconst.AWSVideoS3BuckerName, thumbnailKey)
	}
}
