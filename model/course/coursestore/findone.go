package coursestore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Course, error) {
	var course models.Course
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &course, nil
}
