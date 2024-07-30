package coursestore

import (
	"context"
	"errors"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewCourse(
	ctx context.Context,
	newCourse *models.Course,
) (uint32, error) {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Check if course with the same slug already exists
	var existingCourse models.Course
	if err := tx.Where("slug = ?", newCourse.Slug).First(&existingCourse).Error; err == nil {
		tx.Rollback()
		return 0, errors.New("course with this slug already exists")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return 0, err
	}

	// Create the new course
	if err := tx.Create(newCourse).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return newCourse.Id, nil
}
