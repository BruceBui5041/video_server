package videostore

import (
	"context"
	"video_server/model/video/videomodel"
)

func (s *sqlStore) UpdateVideo(
	ctx context.Context,
	id uint32,
	updateData *videomodel.UpdateVideo,
) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
