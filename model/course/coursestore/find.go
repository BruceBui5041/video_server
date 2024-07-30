package coursestore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindAll(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Course, error) {
	var courses []models.Course
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).Find(&courses).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return courses, nil
}
