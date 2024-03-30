package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetNFTListenersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNFTListenersLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNFTListenersLogic {
	return GetNFTListenersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNFTListenersLogic) GetSyncNftList() (*types.SyncNFTListResponse, error) {
	listeners, err := l.svcCtx.NftListenersModel.FindForSync()
	if err != nil {
		return nil, errorx.NewDefaultError("query failed")
	}
	list := make([]types.SyncNFTListBody, 0)
	for _, listener := range listeners {
		list = append(list, types.SyncNFTListBody{
			Id:              listener.Id,
			Name:            listener.Name.String,
			ContractAddress: listener.ContractAddress,
		})
	}

	return &types.SyncNFTListResponse{
		List: list,
	}, nil
}

func (l *GetNFTListenersLogic) SyncNFT(req types.SyncNFTReq) (*types.SyncNFTResponse, error) {
	err := l.svcCtx.NftListenersModel.SyncNFT(req.ContractAddress, req.TokenId, req.UserAddress, l.svcCtx.NftTokenModel, l.svcCtx.Config.Chain)
	if err != nil {
		return nil, errorx.NewDefaultError("query failed")
	}

	return &types.SyncNFTResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}

func (l *GetNFTListenersLogic) Sync721WithoutTokenId(req types.Sync721WithoutTokenIdReq) (*types.SyncNFTResponse, error) {
	err := l.svcCtx.NftListenersModel.Sync721WithoutTokenId(req.ContractAddress, req.UserAddress, l.svcCtx.NftTokenModel, l.svcCtx.Config.Chain)
	if err != nil {
		return nil, errorx.NewDefaultError("query failed")
	}

	return &types.SyncNFTResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
