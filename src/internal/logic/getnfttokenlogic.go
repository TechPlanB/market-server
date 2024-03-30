package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetNFTTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNFTTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNFTTokenLogic {
	return GetNFTTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNFTTokenLogic) GetNFTTokenIds(req types.NFTTokenIdsReq) (*types.NFTTokenIdsResponse, error) {
	ids, err := l.svcCtx.NftTokenModel.FindTokenIds(req.Owner, req.Contract)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find nft token id list")
	}
	list := make([]int, 0)
	if ids != nil {
		for _, id := range ids {
			list = append(list, id)
		}
	}
	return &types.NFTTokenIdsResponse{
		Code:   200,
		Msg:    "success",
		Result: list,
	}, nil
}

func (l *GetNFTTokenLogic) GetNFTTokenIdsByName(req types.NFTTokenIdsByNameReq) (*types.NFTTokenIdsByNameResponse, error) {
	ids, err := l.svcCtx.NftTokenModel.FindTokenIdsByName(req.Owner, req.Contract, req.Name)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find nft token id list")
	}
	list := make([]int, 0)
	if ids != nil {
		for _, id := range ids {
			list = append(list, id)
		}
	}
	return &types.NFTTokenIdsByNameResponse{
		Code:   200,
		Msg:    "success",
		Result: list,
	}, nil
}
