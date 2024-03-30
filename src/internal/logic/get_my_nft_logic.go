package logic

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"
)

type GetMyNftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyNftLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMyNftLogic {
	return GetMyNftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyNftLogic) GetMyNft(req types.GetMyNftReq) (*types.MyNftResponse, error) {
	myNft, err := l.svcCtx.NftTokenModel.FindInfoByContractAndToken(req.ContractAddress, req.TokenId)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find artwork records")
	}

	myNftBody := &types.MyNftBody{
		Name:        myNft.Name.String,
		Description: myNft.Description.String,
		ImageUrl:    myNft.ImageUrl.String,
	}

	return &types.MyNftResponse{
		Code: 200,
		Msg:  "success",
		Data: *myNftBody,
	}, nil
}
