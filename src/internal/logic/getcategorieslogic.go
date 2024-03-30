package logic

import (
	"context"
	"market/internal/common/errorx"

	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCategoriesLogic {
	return GetCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoriesLogic) GetCategories() (*types.CategoryListResponse, error) {
	categories, err := l.svcCtx.NftCategoriesModel.Find()
	if err != nil {
		return nil, errorx.NewDefaultError("filter failed")
	}
	list := make([]types.CategoryListBody, 0)
	for _, category := range categories {
		list = append(list, types.CategoryListBody{
			Id:   category.Id,
			Code: category.Code,
			Name: category.Name,
		})
	}

	return &types.CategoryListResponse{
		List: list,
	}, nil
}
