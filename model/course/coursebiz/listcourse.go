package coursebiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type CourseStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Course, error)
}

type courseBiz struct {
	courseStore CourseStore
}

func NewCourseBiz(courseStore CourseStore) *courseBiz {
	return &courseBiz{courseStore: courseStore}
}

func (biz *courseBiz) ListCourses(ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Course, error) {
	courses, err := biz.courseStore.FindAll(ctx, conditions, moreInfo...)

	if err != nil {
		return nil, common.ErrCannotListEntity(models.CourseEntityName, err)
	}

	return courses, nil
}
