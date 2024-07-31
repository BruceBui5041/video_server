package videostore

import (
	"context"
	"errors"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewVideo(
	ctx context.Context,
	newVideo *models.Video,
) (uint32, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	var existingVideo models.Video
	if err := tx.Where("slug = ? AND course_id = ?", newVideo.Slug, newVideo.CourseID).First(&existingVideo).Error; err == nil {
		tx.Rollback()
		return 0, errors.New("video with this slug already exists in the course")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Create(newVideo).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return newVideo.Id, nil
}
