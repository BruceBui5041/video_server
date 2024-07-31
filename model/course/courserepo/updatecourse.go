package courserepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
)

type UpdateCourseStore interface {
	Update(
		ctx context.Context,
		id uint32,
		updateData *models.Course,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Course, error)
}

type updateCourseRepo struct {
	updateCourseStore UpdateCourseStore
	categoryStore     CategoryStore
}

func NewUpdateCourseRepo(updateCourseStore UpdateCourseStore, categoryStore CategoryStore) *updateCourseRepo {
	return &updateCourseRepo{updateCourseStore: updateCourseStore, categoryStore: categoryStore}
}

func (repo *updateCourseRepo) UpdateCourse(ctx context.Context, id uint32, input *coursemodel.UpdateCourse) (*models.Course, error) {
	updateData := &models.Course{
		Title:       input.Title,
		Description: input.Description,
		Slug:        input.Slug,
	}

	if input.CategoryID != "" {
		uid, err := common.FromBase58(input.CategoryID)
		if err != nil {
			return nil, err
		}
		category, err := repo.categoryStore.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
		if err != nil {
			return nil, err
		}
		updateData.CategoryID = category.Id
	}

	if err := repo.updateCourseStore.Update(ctx, id, updateData); err != nil {
		return nil, err
	}

	course, err := repo.updateCourseStore.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	return course, nil
}
