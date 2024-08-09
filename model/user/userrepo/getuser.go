package userrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type GetUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type getUserRepo struct {
	store GetUserStore
}

func NewGetUserRepo(store GetUserStore) *getUserRepo {
	return &getUserRepo{store: store}
}

func (repo *getUserRepo) GetUser(ctx context.Context, id uint32) (*models.User, error) {
	user, err := repo.store.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.UserEntityName, err)
	}
	return user, nil
}
