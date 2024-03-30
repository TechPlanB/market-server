package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type NFTDropsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNFTDropsLogic(ctx context.Context, svcCtx *svc.ServiceContext) NFTDropsLogic {
	return NFTDropsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NFTDropsLogic) NFTDrops(req types.NFTDropsReq) (*types.NFTDropsResponse, error) {
	drops, total, err := l.svcCtx.NftDropsModel.Find(req)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch list")
	}

	list := make([]types.NFTDropsBody, 0)
	for _, drop := range drops {
		status := STATUS_ON_SALE
		if drop.StartTime > time.Now().Unix() {
			status = STATUS_NOT_START
		} else {
			if drop.Count <= 0 {
				status = STATUS_SOLD_OUT
			} else {
				if drop.EndTime.Valid && drop.EndTime.Int64 < time.Now().Unix() {
					status = STATUS_ENDED
				}
			}
		}

		list = append(list, types.NFTDropsBody{
			ID:        drop.Id,
			Name:      drop.Name,
			ImageUrl:  drop.ImageUrl,
			Count:     drop.Count,
			Price:     drop.Price,
			StartTime: drop.StartTime,
			EndTime:   drop.EndTime.Int64,
			Status:    status.String(),
		})
	}

	return &types.NFTDropsResponse{
		Code:  200,
		Msg:   "success",
		Total: total,
		List:  list,
	}, nil
}
