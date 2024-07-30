package courserepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
)

type CategoryStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type CreateCourseStore interface {
	CreateNewCourse(
		ctx context.Context,
		newCourse *models.Course,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Course, error)
}

type courseRepo struct {
	courseStore   CreateCourseStore
	categoryStore CategoryStore
}

func NewCourseRepo(courseStore CreateCourseStore, categoryStore CategoryStore) *courseRepo {
	return &courseRepo{courseStore: courseStore, categoryStore: categoryStore}
}

func (repo *courseRepo) CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) (*models.Course, error) {
	uid, err := common.FromBase58(input.CategoryID)
	if err != nil {
		return nil, err
	}

	category, err := repo.categoryStore.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
	if err != nil {
		return nil, err
	}

	// Create new course
	newCourse := &models.Course{
		Title:       input.Title,
		Description: input.Description,
		CreatorID:   input.CreatorID,
		CategoryID:  category.Id,
		Slug:        input.Slug,
	}

	courseId, err := repo.courseStore.CreateNewCourse(ctx, newCourse)
	if err != nil {
		return nil, err
	}

	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"id": courseId})
	if err != nil {
		return nil, err
	}

	return course, nil
}
