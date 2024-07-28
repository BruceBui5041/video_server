package coursebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
)

type CreateCourseStore interface {
	CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) error
}

type createCourseBiz struct {
	createCourseStore CreateCourseStore
}

func NewCreateCourseBiz(createCourseStore CreateCourseStore) *createCourseBiz {
	return &createCourseBiz{createCourseStore: createCourseStore}
}

func (c *createCourseBiz) CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) error {
	// Validate required fields
	if input.Title == "" {
		return errors.New("course title is required")
	}

	if input.CreatorID == 0 {
		return errors.New("creator ID is required")
	}

	if input.CategoryID == 0 {
		return errors.New("category ID is required")
	}

	if input.Slug == "" {
		return errors.New("course slug is required")
	}

	// You might want to add more validation here, such as:
	// - Check if the creator exists
	// - Check if the category exists
	// - Validate the length of the title and description
	// - Check for profanity or inappropriate content

	if len(input.Title) > 255 {
		return errors.New("course title must not exceed 255 characters")
	}

	if len(input.Slug) > 255 {
		return errors.New("course slug must not exceed 255 characters")
	}

	// Call the storage layer to create the course
	err := c.createCourseStore.CreateNewCourse(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.CourseEntityName, err)
	}

	return nil
}
