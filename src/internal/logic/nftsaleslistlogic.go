package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type NFTSalesListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNFTSalesListLogic(ctx context.Context, svcCtx *svc.ServiceContext) NFTSalesListLogic {
	return NFTSalesListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NFTSalesListLogic) NFTSalesList(req types.NFTSalesReq) (*types.NFTSalesResponse, error) {
	sales, count, err := l.svcCtx.NftSalesModel.Find(req)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch list")
	}
	list := make([]types.NFTSalesBody, 0)
	for _, sale := range sales {
		list = append(list, types.NFTSalesBody{
			ID:           sale.Id,
			Name:         sale.Name,
			ImageUrl:     sale.ImageUrl,
			Count:        sale.Count,
			FixedPrice:   sale.FixedPrice,
			HighestPrice: sale.HighestPrice,
			StartTime:    sale.StartTime.Int64,
			EndTime:      sale.EndTime.Int64,
			Status:       sale.Status,
			SaleType:     sale.SaleType,
			TokenId:      sale.TokenId,
			SaleAddress:  sale.SaleAddress,
		})
	}

	return &types.NFTSalesResponse{
		Code:  200,
		Msg:   "success",
		Total: count,
		List:  list,
	}, nil
}
