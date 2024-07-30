package categorybiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

type CategoryStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Category, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type categoryBiz struct {
	categoryStore CategoryStore
}

func NewCategoryBiz(categoryStore CategoryStore) *categoryBiz {
	return &categoryBiz{categoryStore: categoryStore}
}

func (biz *categoryBiz) ListCategories(ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Category, error) {
	categories, err := biz.categoryStore.FindAll(ctx, conditions, moreInfo...)

	if err != nil {
		return nil, common.ErrCannotListEntity(models.CategoryEntityName, err)
	}

	return categories, nil
}
