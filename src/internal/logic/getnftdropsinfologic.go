package logic

import (
	"context"
	"market/internal/common/errorx"
	"time"

	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetNFTDropsInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNFTDropsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNFTDropsInfoLogic {
	return GetNFTDropsInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNFTDropsInfoLogic) GetNFTDropsInfo(req types.NFTDropsInfoReq) (*types.NFTDropsInfoResponse, error) {
	drops, err := l.svcCtx.NftDropsModel.FindInfo(req.ID)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find nft drops records")
	}

	status := STATUS_ON_SALE
	if drops.StartTime > time.Now().Unix() {
		status = STATUS_NOT_START
	} else {
		if drops.Count <= 0 {
			status = STATUS_SOLD_OUT
		} else {
			if drops.EndTime.Valid && drops.EndTime.Int64 < time.Now().Unix() {
				status = STATUS_ENDED
			}
		}
	}

	NFTDropsInfo := &types.NFTDropsInfoBody{
		ID:             drops.Id,
		Name:           drops.Name,
		Description:    drops.Description,
		CreatedAt:      drops.CreatedAt.String(),
		CategoryName:   drops.CategoryName,
		CollectionName: drops.CollectionName,
		ImageUrl:       drops.ImageUrl,
		Count:          drops.Count,
		Price:          drops.Price,
		BuyAddress:     drops.BuyAddress,
		IdInContract:   drops.IdInContract,
		NeedCheck:      drops.NeedCheck,
		ShowLeft:       drops.ShowLeft,
		StartTime:      drops.StartTime,
		EndTime:        drops.EndTime.Int64,
		Status:         status.String(),
		NftAddress:     drops.NftAddress,
	}

	return &types.NFTDropsInfoResponse{
		Code: 200,
		Msg:  "success",
		Data: *NFTDropsInfo,
	}, nil
}
