package coursestore

import (
	"context"
	"errors"
	models "video_server/model"
	"video_server/model/course/coursemodel"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewCourse(
	ctx context.Context,
	input *coursemodel.CreateCourse,
) error {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Check if course with the same title already exists
	var existingCourse models.Course
	if err := tx.Where("slug = ?", input.Slug).First(&existingCourse).Error; err == nil {
		tx.Rollback()
		return errors.New("course with this slug already exists")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	// Create new course
	newCourse := models.Course{
		Title:       input.Title,
		Description: input.Description,
		CreatorID:   input.CreatorID,
		CategoryID:  input.CategoryID,
		Slug:        input.Slug,
	}

	if err := tx.Create(&newCourse).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
