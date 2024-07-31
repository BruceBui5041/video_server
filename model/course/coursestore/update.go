package coursestore

import (
	"context"
	"errors"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	updateData *models.Course,
) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if updateData.Slug != "" {
		var existingCourse models.Course
		if err := tx.Where("slug = ? AND id != ?", updateData.Slug, id).First(&existingCourse).Error; err == nil {
			tx.Rollback()
			return errors.New("course with this slug already exists")
		} else if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&models.Course{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
