package coursebiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type ListCourseStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Course, error)
}

type listCourseBiz struct {
	listCourseStore ListCourseStore
}

func NewCourseBiz(listCourseStore ListCourseStore) *listCourseBiz {
	return &listCourseBiz{listCourseStore: listCourseStore}
}

func (biz *listCourseBiz) ListCourses(ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Course, error) {
	courses, err := biz.listCourseStore.FindAll(ctx, conditions, moreInfo...)

	if err != nil {
		return nil, common.ErrCannotListEntity(models.CourseEntityName, err)
	}

	return courses, nil
}
