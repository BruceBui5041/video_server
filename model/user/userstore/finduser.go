package userstore

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error) {
	db := s.db.Table(models.User{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user models.User

	if err := db.Where(conditions).First(&user).Error; err != nil {
		return nil, common.ErrCannotListEntity(models.CategoryEntityName, err)
	}

	return &user, nil
}
