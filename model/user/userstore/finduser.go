package userstore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound(models.User{}.TableName(), err)
		}
		return nil, common.ErrDB(err)
	}

	return &user, nil
}
