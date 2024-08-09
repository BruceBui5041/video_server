package userbiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type GetUserRepo interface {
	GetUser(ctx context.Context, id uint32) (*models.User, error)
}

type getUserBiz struct {
	repo GetUserRepo
}

func NewGetUserBiz(repo GetUserRepo) *getUserBiz {
	return &getUserBiz{repo: repo}
}

func (biz *getUserBiz) GetUserById(ctx context.Context, id uint32) (*models.User, error) {
	user, err := biz.repo.GetUser(ctx, id)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.UserEntityName, err)
	}
	return user, nil
}
