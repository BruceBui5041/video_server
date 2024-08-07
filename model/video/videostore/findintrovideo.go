package videostore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindIntroVideo(ctx context.Context, courseID uint32) (*models.Video, error) {
	var video models.Video
	if err := s.db.Where("course_id = ? AND intro_video = ?", courseID, true).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, err
	}
	return &video, nil
}
